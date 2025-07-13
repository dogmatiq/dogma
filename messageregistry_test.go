package dogma_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestRegisteredMessageTypeFor(t *testing.T) {
	t.Run("it returns the type that represents T", func(t *testing.T) {
		const id = "2d8ce56f-1983-44e3-a55d-f74d8dcb0adc"
		type T struct{ Command }
		RegisterCommand[T](id)

		mt, ok := RegisteredMessageTypeFor[T]()
		if !ok {
			t.Fatal("expected message type to be registered")
		}

		if got, want := mt.GoType(), reflect.TypeFor[T](); got != want {
			t.Fatalf("unexpected type: got %v, want %v", got, want)
		}
	})

	t.Run("it returns false when T is not in the registry", func(t *testing.T) {
		type T struct{ Command }

		_, ok := RegisteredMessageTypeFor[T]()
		if ok {
			t.Fatal("did not expect message type to be registered")
		}
	})
}

func TestRegisteredMessageTypeByID(t *testing.T) {
	t.Run("it returns the type associated with the normalized ID", func(t *testing.T) {
		const id = "37264B0D-4342-4708-8263-60D82DE78AD1"
		type T struct{ Event }
		RegisterEvent[T](id)

		mt, ok := RegisteredMessageTypeByID(id)
		if !ok {
			t.Fatal("expected message type to be registered")
		}

		if got, want := mt.GoType(), reflect.TypeOf(T{}); got != want {
			t.Fatalf("unexpected type: got %v, want %v", got, want)
		}
	})

	t.Run("it returns false when ID is not registered", func(t *testing.T) {
		const id = "75285aae-f85a-435b-ad36-7471ab169348"

		_, ok := RegisteredMessageTypeByID(id)
		if ok {
			t.Fatal("did not expect message type to be registered")
		}
	})

	t.Run("panics when the ID is invalid", func(t *testing.T) {
		expectPanic(
			t,
			`"<non-uuid>" is not a canonical RFC 4122 UUID: expected 36 characters`,
			func() {
				RegisteredMessageTypeByID("<non-uuid>")
			},
		)
	})
}

func TestRegisteredMessageTypes(t *testing.T) {
	t.Run("yields the registered message types", func(t *testing.T) {
		type T struct{ Command }
		type U struct{ Event }
		type V struct{ Timeout }

		RegisterCommand[T]("b3160ff8-f19a-4f79-b81c-0551c99aeac2")
		RegisterEvent[U]("c3f856ba-0519-4335-ad84-313aa0fedc5e")
		RegisterTimeout[V]("8ab13db4-33ff-4dde-862c-7c94a0477231")

		var yieldedT, yieldedU, yieldedV bool

		for mt := range RegisteredMessageTypes() {
			switch mt.GoType() {
			case reflect.TypeFor[T]():
				yieldedT = true
			case reflect.TypeFor[U]():
				yieldedU = true
			case reflect.TypeFor[V]():
				yieldedV = true
			}
		}

		if !yieldedT {
			t.Fatal("command type was not yielded")
		}

		if !yieldedU {
			t.Fatal("event type was not yielded")
		}

		if !yieldedV {
			t.Fatal("timeout type was not yielded")
		}
	})
}

func TestMessageTypeRegistration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const id = "476497D0-2132-4EC5-90D4-21E2B6E6EB54"
		type T struct{ Command }
		RegisterCommand[T](id)

		mt, ok := RegisteredMessageTypeFor[T]()
		if !ok {
			t.Fatal("message type is not registered")
		}

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
	t.Run("func ID()", func(t *testing.T) {
		t.Run("returns the normalized UUID", func(t *testing.T) {
			const id = "5211A466-010A-4C89-BF36-9A95896BFE2B"
			type T struct{ Command }
			RegisterCommand[T](id)

			mt, ok := RegisteredMessageTypeFor[T]()
			if !ok {
				t.Fatal("message type is not registered")
			}

			if got, want := mt.ID(), strings.ToLower(id); got != want {
				t.Fatalf("unexpected UUID: got %q, want %q", got, want)
			}
		})
	})

	t.Run("func GoType()", func(t *testing.T) {
		t.Run("returns the reflect.Type of the message type", func(t *testing.T) {
			type T struct{ Command }
			RegisterCommand[T]("c1d2e3f4-5678-90ab-cdef-1234567890ab")

			mt, ok := RegisteredMessageTypeFor[T]()
			if !ok {
				t.Fatal("message type is not registered")
			}

			got := mt.GoType()
			want := reflect.TypeFor[T]()

			if got != want {
				t.Fatalf("unexpected message type: got %s, want %s", got, want)
			}
		})
	})

	t.Run("func New()", func(t *testing.T) {
		t.Run("when the type uses non-pointer receivers", func(t *testing.T) {
			const id = "7c5724b3-bce9-413a-9777-94eff973539d"
			type T struct{ Command }
			RegisterCommand[T](id)

			mt, ok := RegisteredMessageTypeFor[T]()
			if !ok {
				t.Fatal("message type is not registered")
			}

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

		t.Run("when the type uses pointer receivers", func(t *testing.T) {
			const id = "7da02018-2a02-44ec-aa1f-b68d66d4887d"
			type T struct {
				messageWithPointerRecievers[CommandValidationScope]
			}

			RegisterCommand[*T](id)

			mt, ok := RegisteredMessageTypeFor[*T]()
			if !ok {
				t.Fatal("message type is not registered")
			}

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
