package service

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/olebedev/emitter"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

var _ interfaces.RegistryService = (*RegistryImpl)(nil)

const registryBucketName = "registry"

type RegistryImpl struct {
	node    interfaces.Node
	logger  *zap.Logger
	streams structs.Map
	subs    structs.Map
}

func (r *RegistryImpl) Node() interfaces.Node {
	return r.node
}

func (r *RegistryImpl) Start() error {
	return nil
}

func (r *RegistryImpl) Stop() error {
	return nil
}

func (r *RegistryImpl) Init() error {
	return utils.CreateBucket(registryBucketName, r.node.Db())
}

func NewRegistry(node interfaces.Node) *RegistryImpl {
	return &RegistryImpl{
		node:    node,
		logger:  node.Logger(),
		streams: structs.NewMap(),
		subs:    structs.NewMap(),
	}
}
func (r *RegistryImpl) Set(sre interfaces.SignedRegistryEntry, trusted bool, receivedFrom net.Peer) error {
	hash := encoding.NewMultihash(sre.PK())
	hashString, err := hash.ToString()
	if err != nil {
		return err
	}
	pid, err := receivedFrom.Id().ToString()
	if err != nil {
		return err
	}
	r.logger.Debug("[registry] set", zap.String("pk", hashString), zap.Uint64("revision", sre.Revision()), zap.String("receivedFrom", pid))

	if !trusted {
		if len(sre.PK()) != 33 {
			return errors.New("Invalid pubkey")
		}
		if int(sre.PK()[0]) != int(types.HashTypeEd25519) {
			return errors.New("Only ed25519 keys are supported")
		}
		if sre.Revision() < 0 || sre.Revision() > 281474976710656 {
			return errors.New("Invalid revision")
		}
		if len(sre.Data()) > types.RegistryMaxDataSize {
			return errors.New("Data too long")
		}

		if !sre.Verify() {
			return errors.New("Invalid signature found")
		}
	}

	existingEntry, err := r.getFromDB(sre.PK())
	if err != nil {
		return err
	}

	if existingEntry != nil {
		if receivedFrom != nil {
			if existingEntry.Revision() == sre.Revision() {
				return nil
			} else if existingEntry.Revision() > sre.Revision() {
				updateMessage := protocol.MarshalSignedRegistryEntry(existingEntry)
				err := receivedFrom.SendMessage(updateMessage)
				if err != nil {
					return err
				}
				return nil
			}
		}

		if existingEntry.Revision() >= sre.Revision() {
			return errors.New("Revision number too low")
		}
	}

	key := encoding.NewMultihash(sre.PK())
	keyString, err := key.ToString()
	if err != nil {
		return err
	}

	eventObj, ok := r.streams.Get(keyString)
	if ok {
		event := eventObj.(*emitter.Emitter)
		go event.Emit("fire", sre)
	}

	err = r.node.Db().Update(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte(registryBucketName))
		err := bucket.Put(sre.PK(), protocol.MarshalSignedRegistryEntry(sre))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = r.BroadcastEntry(sre, receivedFrom)
	if err != nil {
		return err
	}

	return nil
}
func (r *RegistryImpl) BroadcastEntry(sre interfaces.SignedRegistryEntry, receivedFrom net.Peer) error {
	hash := encoding.NewMultihash(sre.PK())
	hashString, err := hash.ToString()
	if err != nil {
		return err
	}
	pid, err := receivedFrom.Id().ToString()
	if err != nil {
		return err
	}
	r.logger.Debug("[registry] broadcastEntry", zap.String("pk", hashString), zap.Uint64("revision", sre.Revision()), zap.String("receivedFrom", pid))
	updateMessage := protocol.MarshalSignedRegistryEntry(sre)

	for _, p := range r.node.Services().P2P().Peers().Values() {
		peer, ok := p.(net.Peer)
		if !ok {
			continue
		}
		if receivedFrom == nil || peer.Id().Equals(receivedFrom.Id()) {
			err := peer.SendMessage(updateMessage)
			if err != nil {
				pid, err := peer.Id().ToString()
				if err != nil {
					return err
				}
				r.logger.Error("Failed to send registry broadcast", zap.String("peer", pid), zap.Error(err))
				return err
			}
		}
	}

	return nil
}
func (r *RegistryImpl) SendRegistryRequest(pk []byte) error {
	query := protocol.NewRegistryQuery(pk)

	request, err := msgpack.Marshal(query)
	if err != nil {
		return err
	}

	// Iterate over peers and send the request
	for _, peerVal := range r.node.Services().P2P().Peers().Values() {
		peer, ok := peerVal.(net.Peer)
		if !ok {
			continue
		}
		err := peer.SendMessage(request)
		if err != nil {
			pid, err := peer.Id().ToString()
			if err != nil {
				return err
			}
			r.logger.Error("Failed to send registry request", zap.String("peer", pid), zap.Error(err))
			return err
		}
	}

	return nil
}
func (r *RegistryImpl) Get(pk []byte) (interfaces.SignedRegistryEntry, error) {
	key := encoding.NewMultihash(pk)
	keyString, err := key.ToString()
	if err != nil {
		return nil, err
	}

	if r.subs.Contains(keyString) {
		r.logger.Debug("[registry] get (cached)", zap.String("key", keyString))
		res, err := r.getFromDB(pk)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}

		err = r.SendRegistryRequest(pk)
		if err != nil {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
		return r.getFromDB(pk)
	}

	err = r.SendRegistryRequest(pk)
	if err != nil {
		return nil, err
	}
	r.subs.Put(keyString, key)
	if !r.streams.Contains(keyString) {
		event := &emitter.Emitter{}
		r.streams.Put(keyString, event)
	}

	res, err := r.getFromDB(pk)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	if res == nil {
		r.logger.Debug("[registry] get (cached)", zap.String("key", keyString))
		for i := 0; i < 200; i++ {
			time.Sleep(10 * time.Millisecond)
			res, err := r.getFromDB(pk)
			if err != nil {
				return nil, err
			}
			if res != nil {
				break
			}
		}
	}

	return nil, nil
}

func (r *RegistryImpl) Listen(pk []byte, cb func(sre interfaces.SignedRegistryEntry)) (func(), error) {
	key, err := encoding.NewMultihash(pk).ToString()
	if err != nil {
		return nil, err
	}

	cbProxy := func(event *emitter.Event) {
		sre, ok := event.Args[0].(interfaces.SignedRegistryEntry)
		if !ok {
			r.logger.Error("Failed to cast event to SignedRegistryEntry")
			return
		}

		cb(sre)
	}

	if !r.streams.Contains(key) {
		em := emitter.New(0)
		r.streams.Put(key, em)
		err := r.SendRegistryRequest(pk)
		if err != nil {
			return nil, err
		}
	}
	streamVal, _ := r.streams.Get(key)
	stream := streamVal.(*emitter.Emitter)
	channel := stream.On("fire", cbProxy)

	return func() {
		stream.Off("fire", channel)
	}, nil
}

func (r *RegistryImpl) getFromDB(pk []byte) (sre interfaces.SignedRegistryEntry, err error) {
	err = r.node.Db().View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte(registryBucketName))
		val := bucket.Get(pk)
		if val != nil {
			entry, err := protocol.UnmarshalSignedRegistryEntry(val)
			if err != nil {
				return err
			}
			sre = entry
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sre, nil
}
