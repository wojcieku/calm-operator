//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Client) DeepCopyInto(out *Client) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Client.
func (in *Client) DeepCopy() *Client {
	if in == nil {
		return nil
	}
	out := new(Client)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyMeasurement) DeepCopyInto(out *LatencyMeasurement) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyMeasurement.
func (in *LatencyMeasurement) DeepCopy() *LatencyMeasurement {
	if in == nil {
		return nil
	}
	out := new(LatencyMeasurement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LatencyMeasurement) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyMeasurementList) DeepCopyInto(out *LatencyMeasurementList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LatencyMeasurement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyMeasurementList.
func (in *LatencyMeasurementList) DeepCopy() *LatencyMeasurementList {
	if in == nil {
		return nil
	}
	out := new(LatencyMeasurementList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LatencyMeasurementList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyMeasurementSpec) DeepCopyInto(out *LatencyMeasurementSpec) {
	*out = *in
	if in.Servers != nil {
		in, out := &in.Servers, &out.Servers
		*out = make([]Server, len(*in))
		copy(*out, *in)
	}
	if in.Clients != nil {
		in, out := &in.Clients, &out.Clients
		*out = make([]Client, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyMeasurementSpec.
func (in *LatencyMeasurementSpec) DeepCopy() *LatencyMeasurementSpec {
	if in == nil {
		return nil
	}
	out := new(LatencyMeasurementSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyMeasurementStatus) DeepCopyInto(out *LatencyMeasurementStatus) {
	*out = *in
	out.State = in.State
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyMeasurementStatus.
func (in *LatencyMeasurementStatus) DeepCopy() *LatencyMeasurementStatus {
	if in == nil {
		return nil
	}
	out := new(LatencyMeasurementStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Server) DeepCopyInto(out *Server) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Server.
func (in *Server) DeepCopy() *Server {
	if in == nil {
		return nil
	}
	out := new(Server)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *State) DeepCopyInto(out *State) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new State.
func (in *State) DeepCopy() *State {
	if in == nil {
		return nil
	}
	out := new(State)
	in.DeepCopyInto(out)
	return out
}
