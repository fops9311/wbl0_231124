package storage

type Storage interface {
	StorageGet(id string) (interface{}, bool)
	StorageGetAll() []interface{}
	StorageUpdate(id string, val interface{})
	StorageGetKeys() []string
	StorageInit()
}

type AnyStorage[T any] struct {
	StorageGet     func(id string) (T, bool)
	StorageGetAll  func() []T
	StorageUpdate  func(id string, val T)
	StorageGetKeys func() []string
}

func NewStorage[T any](s Storage) *AnyStorage[T] {
	s.StorageInit()
	res := &AnyStorage[T]{
		StorageGet: func(id string) (T, bool) {
			a, b := s.StorageGet(id)
			t, ok := a.(T)
			if !ok {
				var t2 T
				return t2, false
			}
			return t, b
		},
		StorageGetAll: func() []T {
			sls := s.StorageGetAll()
			var res []T = make([]T, len(sls))
			cnt := 0
			for _, v := range sls {
				t, ok := v.(T)
				if !ok {
					continue
				}
				res[cnt] = t
				cnt++
			}
			res = res[:cnt]
			return res
		},
		StorageUpdate: func(id string, val T) {
			s.StorageUpdate(id, val)
		},
		StorageGetKeys: func() []string {
			return s.StorageGetKeys()
		},
	}
	return res
}
