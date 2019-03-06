/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package output

// Outputer is a destination for PCAPs. Concrete implementations
// can be 'save to file', 'view in wireshark', etc.
type Outputer interface {
	Write(p []byte) (n int, err error)
	Close()
}
