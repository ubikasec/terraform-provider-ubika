package v1beta

// Status is an Unknownable
//func (e *AssetStatusResourceModel) SetUnknown(_ context.Context, state bool) error {
//	e.unknown = state
//	return nil
//}

//func (e *AssetStatusResourceModel) SetValue(ctx context.Context, v interface{}) error {
//	tflog.Info(ctx, "should set value", map[string]interface{}{"value": fmt.Sprintf("%#v", v)})
//	//		if vv, ok := v.(*AssetStatusResourceModel); ok {
//	//			e = vv
//	//		}
//	return fmt.Errorf("SetValue on status detected")
//}
//func (e *AssetStatusResourceModel) GetUnknown(context.Context) bool { return e.unknown }

//func (e AssetStatusResourceModel) GetValue(ctx context.Context) interface{} {
//	return map[string]tftypes.Value{
//		"service_address": tftypes.NewValue(tftypes.String, e.ServiceAddress.ValueString()),
//		// "state":           tftypes.List[tftypes.String],
//	}
//}
