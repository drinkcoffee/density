package main

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "sort"
)

const NUM_NODES = 100
const DHT_SPAN = 16
const WORD_SIZE = 32  // 32 bytes



func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println ("start by generating", NUM_NODES, "random numbers")
    fmt.Println()

    var existingGatewayIds[NUM_NODES] *big.Int

    max := new(big.Int)
    max.Exp(big.NewInt(2), big.NewInt(256), nil).Sub(max, big.NewInt(1))
    
    for i := 0; i < NUM_NODES; i++ {
        //Generate cryptographically strong pseudo-random between 0 - max
        n, err := rand.Int(rand.Reader, max)
        existingGatewayIds[i] = n
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        fmt.Printf("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
    }
    fmt.Println()


    ids := []string{ } 
    for i := 0; i < NUM_NODES; i++ {
        var z string
        z = existingGatewayIds[i].Text(10)
        ids = append(ids, z)
    }


    sort.Strings(ids)

   
    fmt.Println("Sorted list of Gateway Ids")
    for i := range ids {
        fmt.Printf("%v", ids[i])
        fmt.Println()
    }

//     fmt.Println("Sorted list of Gateway Ids")
//     for i := 0; i < NUM_NODES; i++ {
// //        fmt.Printf("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//         fmt.Println("Existing Gateway Id: 0x%x\n", existingGatewayIds[i])
//     }









}
