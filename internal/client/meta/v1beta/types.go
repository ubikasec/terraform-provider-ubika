package metav1beta

// Namespaced object

const collectionType = "ns"

func (m *ObjectMeta) GetCollectionName() string {
	return m.GetNamespace()
}

func (m *ObjectMeta) GetCollectionType() string {
	return collectionType
}

func (o *DeleteOptions) GetCollectionName() string {
	return o.GetNamespace()
}

func (o *GetOptions) GetCollectionName() string {
	return o.GetNamespace()
}

func (o *ListOptions) GetCollectionName() string {
	return o.GetNamespace()
}

func (o *PatchOptions) GetCollectionName() string {
	return o.GetNamespace()
}

func (o *WatchOptions) GetCollectionName() string {
	return o.GetNamespace()
}
