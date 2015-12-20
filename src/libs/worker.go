/**
Copyright 2015 andrew@yasinsky.com

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
package nodder

//worker instantiates
//worker tries to become master,
//if not joins a swarm of bees and retries later
//masters receive and dole out wokloads
//bees do the work return data and die

type Worker struct {
}
