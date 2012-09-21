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
	//	subscriptions map[string]*Subscription
	conn PubSubConnection
}

func NewPubSubClient() (PubSubClient, Error) {
	spec := DefaultSpec().Protocol(REDIS_PUBSUB)
	return NewPubSubClientWithSpec(spec)
}

func NewPubSubClientWithSpec(spec *ConnectionSpec) (PubSubClient, Error) {
	c := new(pubsubClient)
	var err Error
	c.conn, err = NewPubSubConnection(spec)
	if err != nil {
		return nil, err
	}

	//	c.messages = make(chan []byte, spec.rspChanCap)
	//	c.subscriptions = make(map[string]*Subscription)

	return c, nil
}

func (c *pubsubClient) Messages(topic string) PubSubChannel {
	if s := c.conn.Subscriptions()[topic]; s != nil && s.IsActive {
		return s.Channel
	}
	return nil
}

func (c *pubsubClient) Subscriptions() []string {
	topics := make([]string, 0)
	for topic, s := range c.conn.Subscriptions() {
		if s.IsActive {
			topics = append(topics, topic)
		}
	}
	return topics
}

// REVU - why not async semantics?
func (c *pubsubClient) Subscribe(topic string, otherTopics ...string) (err Error) {
	args := appendAndConvert(topic, otherTopics...)
	//	var ok bool
	_, err = c.conn.ServiceRequest(&SUBSCRIBE, args)
	//	if err == nil {
	//		err = NewError(REDIS_ERR, "Subscribe() NOT IMPLEMENTED")
	//	}
	return
}

// REVU - why not async semantics?
func (c *pubsubClient) Unsubscribe(topics ...string) (err Error) {
	if topics == nil {
		topics = c.Subscriptions()
	}
	var otherTopics []string = nil
	if len(topics) > 1 {
		otherTopics = topics[1:]
	}
	args := appendAndConvert(topics[0], otherTopics...)
	//	var ok bool
	_, err = c.conn.ServiceRequest(&UNSUBSCRIBE, args)
	//	if err == nil {
	//		err = NewError(REDIS_ERR, "Subscribe() NOT IMPLEMENTED")
	//	}
	return
}

// REVU - why not async semantics?
func (c *pubsubClient) Quit() Error {
	return NewError(REDIS_ERR, "Quit() NOT IMPLEMENTED")
}
