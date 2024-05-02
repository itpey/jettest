// Copyright 2024 itpey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import "github.com/urfave/cli/v2"

const (
	appNameArt = `_______________________________________________________
______  /__  ____/__  __/__  __/__  ____/_  ___/__  __/
___ _  /__  __/  __  /  __  /  __  __/  _____ \__  /   
/ /_/ / _  /___  _  /   _  /   _  /___  ____/ /_  /    
\____/  /_____/  /_/    /_/    /_____/  /____/ /_/     ` + "\n\n"
	appVersion     = "0.1.0"
	appDescription = `JetTest is a versatile command-line application that allows you to execute API tests defined in YAML configuration files.`
	appCopyright   = "Apache-2.0 license\nFor more information, visit the GitHub repository: https://github.com/itpey/jettest"
)

var (
	appAuthors = []*cli.Author{{Name: "itpey", Email: "itpey@github.com"}}
)
