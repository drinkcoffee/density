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
    "math/big"
    "sort"
    "github.com/drinkcoffee/density/internal/density"
)

const numNodes = 100
const dhtSpan = 16
const wordSize = 32  // 32 bytes



func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println()

    fmt.Println ("Start by generating", numNodes, "random numbers")
    var gatewayIds = density.NewGatewayIds(wordSize)
    for i := 0; i < numNodes; i++ {
        gatewayIds.AddRandom()
    }
    fmt.Println()

    fmt.Println ("Sort them")
    sort.Sort(gatewayIds);
    gatewayIds.PrintAll()


    var numRange = new(big.Int)
    numRange.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)
    
    var temp1 = new(big.Int)
    temp1.Div(numRange, big.NewInt(numNodes))
    var windowSize = new(big.Int)
    windowSize.Mul(temp1, big.NewInt(dhtSpan))

    fmt.Printf("Window Size: 0x%x\n", &windowSize)




    // ids := []string{ } 
    // for i := 0; i < numNodes; i++ {
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
//     for i := 0; i < numNodes; i++ {
// //        fmt.Printf("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//         fmt.Println("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//     }









}
