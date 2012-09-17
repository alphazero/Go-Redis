//   Copyright 2009-2012 Joubin Houshyar
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package redis

import ()

// -----------------------------------------------------------------------------
// pubsubClient - supports PubSubClient interface
// -----------------------------------------------------------------------------

type pubsubClient struct {
	//	messages      chan []byte
	subscriptions map[string]*subscriptionInfo
	conn          PubSubConnection
}

type subscriptionInfo struct {
	channel chan []byte
}

func NewPubSubClient() (PubSubClient, Error) {
	spec := DefaultSpec().Protocol(REDIS_PUBSUB)
	return NewPubSubClientWithSpec(spec)
}

func NewPubSubClientWithSpec(spec *ConnectionSpec) (PubSubClient, Error) {
	c := new(pubsubClient)
	spec.Protocol(REDIS_PUBSUB) // must be so set it regardless
	var err Error
	c.conn, err = NewAsynchConnection(spec)
	if err != nil {
		return nil, err
	}

	//	c.messages = make(chan []byte, spec.rspChanCap)
	c.subscriptions = make(map[string]*subscriptionInfo)

	return c, nil
}

func (psc *pubsubClient) Messages(id string) PubSubChannel {
	if s := psc.subscriptions[id]; s != nil {
		return s.channel
	}
	return nil
}

func (psc *pubsubClient) Subscriptions() []string {
	ids := make([]string, 0)
	for id, _ := range psc.subscriptions {
		ids = append(ids, id)
	}
	return ids
}

func (psc *pubsubClient) Subscribe(channel string, otherChannels ...string) (ch PubSubChannel, subscriptionCount int, err Error) {
	//	args := appendAndConvert(channel, otherChannels...)

	err = NewError(REDIS_ERR, "Subscribe() NOT IMPLEMENTED")
	return
}

func (psc *pubsubClient) Unsubscribe(channels ...string) (subscriptionCount int, err Error) {
	err = NewError(REDIS_ERR, "Unsubscribe() NOT IMPLEMENTED")
	return
}

func (psc *pubsubClient) Quit() Error {
	return NewError(REDIS_ERR, "Quit() NOT IMPLEMENTED")
}
