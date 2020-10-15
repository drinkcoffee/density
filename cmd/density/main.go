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
    "github.com/drinkcoffee/density/internal/density"
)

const initialNumIdentifiers = 64
const dhtSpan = 16
const wordSize = 32  // 32 bytes


func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println()

    var max big.Int
    max.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)

    fmt.Println ("Start by generating", initialNumIdentifiers, "random numbers")
    randomIdentifiers := []big.Int{ } 
    for i := 0; i < initialNumIdentifiers; i++ {
        randomIdentifiers = append(randomIdentifiers, *density.CreateRandomIdentifier(&max))
    }
    ids := density.NewSortedArrayBigInts()
    ids.AddMany(randomIdentifiers)
    fmt.Println(" Generated ids")
    ids.PrintAll()
    fmt.Println()


    var gatewayIds = density.NewValueRing(wordSize, dhtSpan, ids)
    fmt.Println()

    for i := 0; i < 10; i++ {
        ideal := gatewayIds.FindBestLocationForValue()
        fmt.Printf("Ideal:     %#x\n", &ideal)
    
        ids.Add(ideal)
    }
}