// To parse and unparse this JSON data, add this code to your project and do:
//
//    call, err := UnmarshalCall(bytes)
//    bytes, err = call.Marshal()

package main

import (
	"bytes"
	"encoding/gob"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"

	"github.com/Necroforger/dgrouter/exrouter"
)

// CallID represents an incomplete call.
type CallID struct {
	From string   `json:"From"`
	To   []string `json:"To"`
}

// Call represents an incomplete call.
type Call struct {
	From *discordgo.Channel   `json:"From"`
	To   []*discordgo.Channel `json:"To"`
}

func encodeString(s string) (b []byte, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewEncoder(buffer).Encode(s)
	b = buffer.Bytes()
	return
}

func decodeString(b []byte) (s string, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewDecoder(buffer).Decode(&s)
	return
}

func encodeStrings(s []string) (b []byte, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewEncoder(buffer).Encode(s)
	b = buffer.Bytes()
	return
}

func decodeStrings(b []byte) (s []string, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewDecoder(buffer).Decode(&s)
	return
}

func toChannels(ctx *exrouter.Context, ids []string, onErr func(string, error)) (channels []*discordgo.Channel) {
	for _, id := range ids {
		channel, err := ctx.Channel(id)
		if err != nil {
			onErr(id, err)
			continue
		}
		channels = append(channels, channel)
	}
	return
}

func insertCall(ctx *exrouter.Context, from string, to []string, txn *badger.Txn) (err error) {

	// String to bytes
	bFrom, err := encodeString(from)
	if err != nil {
		logctx(ctx, "Couldn't turn the Channel ID into bytes.", nil)
		return err
	}

	// Strings to bytes
	bTo, err := encodeStrings(to)
	if err != nil {
		logctx(ctx, "Couldn't turn Channel IDs into bytes.", nil)
		return err
	}

	// Set
	err = txn.Set(bFrom, bTo)
	if err != nil {
		logctx(ctx, "Couldn't save this channel's calls.", nil)
		return err
	}

	return err
}

func appendCall(ctx *exrouter.Context, from string, to string, txn *badger.Txn, item *badger.Item) (duplicate bool, err error) {
	err = item.Value(func(val []byte) (err error) {

		// Get called channels
		tos, err := decodeStrings(val)
		if err != nil {
			logctx(ctx, "Couldn't decode Channel IDs.", err)
			return
		}

		// Check for duplicates
		duplicate = csia(tos, to)
		if duplicate {
			// Don't update
		} else {

			// Add this one
			tos = append(tos, to)

			// Save them
			return insertCall(ctx, from, tos, txn)
		}
		return
	})
	return
}

func selectCall(ctx *exrouter.Context, channel *discordgo.Channel) (call *CallID, err error) {
	bid, err := encodeString(channel.ID)
	if err != nil {
		logctx(ctx, "Couldn't encode a Channel ID.", nil)
		return
	}

	err = db.View(func(txn *badger.Txn) (err error) {

		// Get a channel's calls
		item, err := txn.Get(bid)
		if err != nil {
			logctx(ctx, "Couldn't get a call.", nil)
			return
		}

		item.Value(func(val []byte) (err error) {

			// Get called channels
			tos, err := decodeStrings(val)
			if err != nil {
				logctx(ctx, "Couldn't decode Channel IDs.", nil)
				return
			}

			// Return the Call IDs
			call = &CallID{
				From: channel.ID,
				To:   tos,
			}

			return
		})
		return
	})
	return
}

func (callid CallID) toChannels(ctx *exrouter.Context, onErr func(string, error)) (call *Call) {
	from, err := ctx.Channel(callid.From)
	if err != nil {
		onErr(callid.From, err)
	}
	call.From = from
	call.To = toChannels(ctx, callid.To, onErr)
	return
}
