/*
Copyright 2018 codestation

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

package stores

// Storer represents the methods to store/retrieve a backup from another location
type Storer interface {
	Store(filepath string, filename string) error
	Retrieve(s3path string) (string, error)
	RemoveOlderBackups(keep int) error
	FindLatestBackup() (string, error)
	Close()
}
