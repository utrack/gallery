package hub

import (
	"encoding/json"
	"errors"
	"github.com/pquerna/ffjson/ffjson"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/utrack/gallery/messages"
	"strconv"
	"testing"
	"time"
)

func TestHub(t *testing.T) {
	Convey("With sample hub", t, func() {
		// init dependencies' mocks
		lister := &listerMock{}
		lister.retList = messages.FilesInfo{
			messages.FileInfo{Filename: "somefilename"},
		}
		notifier := &notifierMock{retChan: make(chan messages.FileChangeNotification, 10)}

		// init the hub
		h := NewHub(lister, notifier).(*hub)
		So(h.conns, ShouldHaveLength, 0)

		Convey("Shouldn't register the conn if error was returned", func() {
			lister.retErr = errors.New("test error")

			conn := &connMock{
				rcvChan:    make(chan json.RawMessage, 5),
				disconChan: make(chan error, 5),
			}

			err := h.Accept(conn)
			So(err, ShouldNotBeNil)

			Convey("Connmap's length should be 0", func() {
				So(h.conns, ShouldHaveLength, 0)
			})
		})

		Convey("Should register conns successfully", func() {
			conn := &connMock{
				rcvChan:    make(chan json.RawMessage, 5),
				disconChan: make(chan error, 5),
			}

			err := h.Accept(conn)
			So(err, ShouldBeNil)

			Convey("Conn should be in the conns' map", func() {
				So(h.conns, ShouldHaveLength, 1)
			})
			Convey("lastConnId should be bumped", func() {
				So(h.lastConnId, ShouldEqual, uint64(1))
			})

			Convey("Should push the Lister's contents on acception", func() {
				var got []byte
				// try to retrieve the msg
				select {
				case <-time.After(time.Second):
				case got = <-conn.rcvChan:
				}
				So(len(got), ShouldNotEqual, 0)

				// unmarshal to the resembling struct
				var gotArray messages.FilesInfo
				err := ffjson.Unmarshal(got, &gotArray)
				So(err, ShouldBeNil)
				So(gotArray, ShouldResemble, lister.retList)

			})

			Convey("With flushed seed's contents", func() {
				select {
				case <-time.After(time.Second):
					So(0, ShouldEqual, 1)
				case <-conn.rcvChan:
				}

				Convey("Should push the notification successfully", func() {
					sent := messages.FileChangeNotification{
						Filename: "test filename",
						Action:   messages.ChangeModification,
					}
					notifier.retChan <- sent

					var gotJson json.RawMessage

					select {
					case <-time.After(time.Second):
					case gotJson = <-conn.rcvChan:
					}
					So(len(gotJson), ShouldNotEqual, 0)

					var got messages.FileChangeNotification
					err := ffjson.Unmarshal(gotJson, &got)
					So(err, ShouldBeNil)
					So(got, ShouldResemble, sent)
				})
			})

			Convey("Disconnection", func() {
				conn.disconChan <- errors.New("test discon")
				<-time.After(time.Second / 2)

				h.connMu.RLock()
				So(h.conns, ShouldHaveLength, 0)
				h.connMu.RUnlock()
			})

			Convey("With second conn", func() {
				conn2 := &connMock{
					rcvChan:    make(chan json.RawMessage, 5),
					disconChan: make(chan error, 5),
				}

				err := h.Accept(conn2)
				So(err, ShouldBeNil)

				Convey("Conn should be in the conns' map", func() {
					h.connMu.RLock()
					So(h.conns, ShouldHaveLength, 2)
					h.connMu.RUnlock()
				})

				Convey("With flushed seed contents", func() {
					<-conn.rcvChan
					<-conn2.rcvChan

					Convey("Notifications should be sent to everyone", func() {
						sent := messages.FileChangeNotification{
							Filename: "test filename",
							Action:   messages.ChangeModification,
						}
						notifier.retChan <- sent

						for num, c := range []*connMock{conn, conn2} {

							Convey("Conn "+strconv.Itoa(num), func() {
								var gotJson json.RawMessage

								select {
								case <-time.After(time.Second):
								case gotJson = <-c.rcvChan:
								}
								So(len(gotJson), ShouldNotEqual, 0)

								var got messages.FileChangeNotification
								err := ffjson.Unmarshal(gotJson, &got)
								So(err, ShouldBeNil)
								So(got, ShouldResemble, sent)
							})
						}
					})
				})
			})
		})

	})
}
