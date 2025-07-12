package dogma_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestMessageTypeRegistration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const id = "476497D0-2132-4EC5-90D4-21E2B6E6EB54"
		type T struct{ Command }
		RegisterCommand[T](id)

		mt := RegisteredMessageTypeFor[T]()

		t.Run("normalizes the UUID", func(t *testing.T) {
			if got, want := mt.ID(), strings.ToLower(id); got != want {
				t.Fatalf("non-normalized UUID: got %q, want %q", got, want)
			}
		})
	})

	t.Run("failure conditions", func(t *testing.T) {
		cases := []struct {
			Name         string
			CommandError string
			CommandFunc  func()
			EventError   string
			EventFunc    func()
			TimeoutError string
			TimeoutFunc  func()
		}{
			{
				"duplicate registration",
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered`,
				func() {
					type T struct{ Command }
					RegisterCommand[T]("658e7da4-ffa9-4d11-b796-29346ec6f586")
					RegisterCommand[T]("658e7da4-ffa9-4d11-b796-29346ec6f586")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered`,
				func() {
					type T struct{ Event }
					RegisterEvent[T]("746979d2-2ce3-4b2f-a4d0-d34b069ba31b")
					RegisterEvent[T]("746979d2-2ce3-4b2f-a4d0-d34b069ba31b")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered`,
				func() {
					type T struct{ Timeout }
					RegisterTimeout[T]("6dfb9656-92b4-4cef-b51d-f808bf6403b2")
					RegisterTimeout[T]("6dfb9656-92b4-4cef-b51d-f808bf6403b2")
				},
			},
			{
				"conflicting registration (same type)",
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered as "f2a5c633-fbc1-4230-937d-057e3d141d4f"`,
				func() {
					type T struct{ Command }
					RegisterCommand[T]("f2a5c633-fbc1-4230-937d-057e3d141d4f")
					RegisterCommand[T]("6b529039-03fe-4fde-8986-42892f39d93e")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered as "cf7a6893-da73-4638-9438-c59d32b2e087"`,
				func() {
					type T struct{ Event }
					RegisterEvent[T]("cf7a6893-da73-4638-9438-c59d32b2e087")
					RegisterEvent[T]("ffd82f55-a6e3-4c02-b43c-491a5abc1b46")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: it is already registered as "2e4c3894-f76b-4ccb-a817-4422df36138e"`,
				func() {
					type T struct{ Timeout }
					RegisterTimeout[T]("2e4c3894-f76b-4ccb-a817-4422df36138e")
					RegisterTimeout[T]("3de60928-e9be-4a62-9fd8-60734f53cde5")
				},
			},
			{
				"conflicting registration (same ID)",
				`cannot register github.com/dogmatiq/dogma_test.U: "b8a9ee69-e6f4-41ae-8fba-524505a2aaba" is already associated with github.com/dogmatiq/dogma_test.T`,
				func() {
					type T struct{ Command }
					type U struct{ Command }
					RegisterCommand[T]("b8a9ee69-e6f4-41ae-8fba-524505a2aaba")
					RegisterCommand[U]("b8a9ee69-e6f4-41ae-8fba-524505a2aaba")
				},
				`cannot register github.com/dogmatiq/dogma_test.U: "21b3ca92-2667-4f0e-b175-7f30f015f968" is already associated with github.com/dogmatiq/dogma_test.T`,
				func() {
					type T struct{ Event }
					type U struct{ Event }
					RegisterEvent[T]("21b3ca92-2667-4f0e-b175-7f30f015f968")
					RegisterEvent[U]("21b3ca92-2667-4f0e-b175-7f30f015f968")
				},
				`cannot register github.com/dogmatiq/dogma_test.U: "66c69a42-ea81-4ca9-8587-bf88e8abaf34" is already associated with github.com/dogmatiq/dogma_test.T`,
				func() {
					type T struct{ Timeout }
					type U struct{ Timeout }
					RegisterTimeout[T]("66c69a42-ea81-4ca9-8587-bf88e8abaf34")
					RegisterTimeout[U]("66c69a42-ea81-4ca9-8587-bf88e8abaf34")
				},
			},
			{
				"interface type",
				`cannot register github.com/dogmatiq/dogma_test.T: message type is an interface, expected a concrete type`,
				func() {
					type T interface{ Command }
					RegisterCommand[T]("4a43823c-66b7-43d5-a3ac-e7247c83ddd0")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: message type is an interface, expected a concrete type`,
				func() {
					type T interface{ Event }
					RegisterEvent[T]("2d5b4609-825d-45be-ac69-23a9dcbf1bff")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: message type is an interface, expected a concrete type`,
				func() {
					type T interface{ Timeout }
					RegisterTimeout[T]("07799033-587c-4e5f-bc3d-7ad1797e7ff6")
				},
			},
			{
				"invalid UUID",
				`cannot register github.com/dogmatiq/dogma_test.T: "<non-uuid>" is not a canonical RFC 4122 UUID: expected 36 characters`,
				func() {
					type T struct{ Command }
					RegisterCommand[T]("<non-uuid>")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: "<non-uuid>" is not a canonical RFC 4122 UUID: expected 36 characters`,
				func() {
					type T struct{ Event }
					RegisterEvent[T]("<non-uuid>")
				},
				`cannot register github.com/dogmatiq/dogma_test.T: "<non-uuid>" is not a canonical RFC 4122 UUID: expected 36 characters`,
				func() {
					type T struct{ Timeout }
					RegisterTimeout[T]("<non-uuid>")
				},
			},
			{
				"pointer to type that uses non-pointer receivers",
				`cannot register *github.com/dogmatiq/dogma_test.T: message type uses non-pointer receivers, use github.com/dogmatiq/dogma_test.T (non-pointer) instead`,
				func() {
					type T struct{ Command }
					RegisterCommand[*T]("fda0112d-29ce-4533-b03f-5b9dcfc3907a")
				},
				`cannot register *github.com/dogmatiq/dogma_test.T: message type uses non-pointer receivers, use github.com/dogmatiq/dogma_test.T (non-pointer) instead`,
				func() {
					type T struct{ Event }
					RegisterEvent[*T]("5dad0988-498e-47e1-989e-b3014612c9d1")
				},
				`cannot register *github.com/dogmatiq/dogma_test.T: message type uses non-pointer receivers, use github.com/dogmatiq/dogma_test.T (non-pointer) instead`,
				func() {
					type T struct{ Timeout }
					RegisterTimeout[*T]("b94ffd35-3434-4412-b46c-3ace7235595b")
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				t.Run("command kind", func(t *testing.T) {
					expectPanic(t, c.CommandError, c.CommandFunc)
				})

				t.Run("event kind", func(t *testing.T) {
					expectPanic(t, c.EventError, c.EventFunc)
				})

				t.Run("timeout kind", func(t *testing.T) {
					expectPanic(t, c.TimeoutError, c.TimeoutFunc)
				})
			})
		}
	})
}

func TestRegisteredMessageType(t *testing.T) {
	t.Run("non-pointer implementation", func(t *testing.T) {
		const id = "7c5724b3-bce9-413a-9777-94eff973539d"
		type T struct{ Command }
		RegisterCommand[T](id)

		mt := RegisteredMessageTypeFor[T]()

		t.Run("func New()", func(t *testing.T) {
			t.Run("it returns a zero-value", func(t *testing.T) {
				m := mt.New()

				if _, ok := m.(T); !ok {
					t.Fatalf("unexpected type: got %T, want %T", m, T{})
				}

				if v := reflect.ValueOf(m); !v.IsZero() {
					t.Fatalf("expected zero-value message, got %v", v)
				}
			})
		})
	})

	t.Run("pointer implementation", func(t *testing.T) {
		const id = "7da02018-2a02-44ec-aa1f-b68d66d4887d"
		type T struct {
			messageWithPointerRecievers[CommandValidationScope]
		}

		RegisterCommand[*T](id)

		mt := RegisteredMessageTypeFor[*T]()

		t.Run("func New()", func(t *testing.T) {
			t.Run("it returns a pointer to a zero-value", func(t *testing.T) {
				m := mt.New()

				p, ok := m.(*T)
				if !ok {
					t.Fatalf("unexpected type: got %T, want %T", m, (*T)(nil))
				}

				if p == nil {
					t.Fatal("expected non-nil pointer")
				}

				if v := reflect.ValueOf(*p); !v.IsZero() {
					t.Fatalf("expected pointer to zero-value message: got %v", v)
				}
			})
		})
	})
}
