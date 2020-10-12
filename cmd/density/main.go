/**
 * Copyright 2020 Peter Robinson
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *   http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
    "fmt"
    "sort"
    "github.com/drinkcoffee/density/internal/density"
)

const NUM_NODES = 100
const DHT_SPAN = 16
const WORD_SIZE = 32  // 32 bytes



func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println ("start by generating", NUM_NODES, "random numbers")
    fmt.Println()

    var gatewayIds = density.NewGatewayIds(WORD_SIZE)

    
    for i := 0; i < NUM_NODES; i++ {
        gatewayIds.AddRandom()
    }
    fmt.Println()

    sort.Sort(gatewayIds);

    // ids := []string{ } 
    // for i := 0; i < NUM_NODES; i++ {
    //     var z string
    //     z = existingGatewayIds[i].Text(10)
    //     ids = append(ids, z)
    // }


    // sort.Strings(ids)

   
    // fmt.Println("Sorted list of Gateway Ids")
    // for i := range ids {
    //     fmt.Printf("%v", ids[i])
    //     fmt.Println()
    // }



//     fmt.Println("Sorted list of Gateway Ids")
//     for i := 0; i < NUM_NODES; i++ {
// //        fmt.Printf("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//         fmt.Println("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//     }









}
