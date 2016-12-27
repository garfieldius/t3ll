package log

/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"fmt"
	"os"
	"time"
)

const TimeLayout string = "2006-01-02 15:04:05"

func now() string {
	return time.Now().Format(TimeLayout)
}

func format(severity, msg string, a ...interface{}) string {
	full := fmt.Sprintf("%s [%s] %s\n", now(), severity, msg)
	return fmt.Sprintf(full, a...)
}

func Fatal(msg string, a ...interface{}) {
	panic(format("ERROR", msg, a...))
}

func Err(msg string, a ...interface{}) {
	fmt.Fprint(os.Stderr, format("ERROR", msg, a...))
}
