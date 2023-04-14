package runtime

import "google.golang.org/protobuf/proto"

type GroupVersion struct {
	Group   string
	Version string
}

func (gv GroupVersion) String() string {
	return gv.Group + "/" + gv.Version
}

func (gv GroupVersion) WithKind(kind string) GroupVersionKind {
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

func (gvk GroupVersionKind) String() string {
	return gvk.Group + "/" + gvk.Version + ", Kind=" + gvk.Kind
}

func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

type Object interface {
	proto.Message
	GroupVersionKind() GroupVersionKind
	Validate() error
}

type NamedObject interface {
	Object
	GetName() string
	SetName(string)
	GetVersion() int64
	SetVersion(int64)
}
