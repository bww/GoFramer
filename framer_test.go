// 
// Copyright (c) 2014 Brian William Wolter, All rights reserved.
// Go Framer
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package framer

import (
  "fmt"
  "bytes"
  "testing"
)

func TestReadWrite(t *testing.T) {
  message := "Morbi in faucibus augue, quis tincidunt enim. Aenean accumsan in purus at viverra. Quisque rutrum odio quis est varius fermentum. In mattis diam at dui sodales mollis. Integer placerat vel quam eget volutpat. Ut eu velit sodales dui suscipit egestas. Nulla vel aliquet enim. Nunc nec lacus ac mauris auctor facilisis. Proin lacinia, felis ut suscipit feugiat, augue justo placerat ipsum, vitae sollicitudin tellus quam ut erat. Suspendisse ultrices sapien et quam cursus malesuada. Suspendisse potenti. Fusce eu ligula ac sem tincidunt vulputate. Duis cursus sem vitae aliquet tristique. Donec vitae sagittis odio. Vivamus facilisis, tortor porta auctor tempus, erat quam molestie nisi, eget vehicula nunc justo vitae nisl. Integer scelerisque at urna at ultricies. Aliquam varius velit ut dolor feugiat lobortis. Donec enim enim, porttitor in viverra vel, varius quis neque. Quisque maximus ullamcorper nunc. Vestibulum varius tempus lectus, sed scelerisque mauris sollicitudin non. Quisque ac tellus in lectus luctus malesuada. Curabitur ex dolor, consectetur at dictum vitae, mattis sit amet dui. Donec posuere nulla porttitor, lacinia orci sed, gravida risus. Pellentesque condimentum est purus, nec varius neque aliquam a. Mauris ultrices pretium odio ut ultricies. Etiam id pellentesque libero. Mauris euismod non neque in rhoncus. Phasellus non laoreet urna. Sed vulputate posuere libero. Integer a accumsan mauris. Vivamus nec maximus est. Fusce bibendum lorem sed dapibus finibus. Quisque dapibus tortor euismod bibendum rutrum. Donec iaculis, est in feugiat accumsan, eros magna tincidunt augue, et imperdiet risus libero et risus. Aliquam erat volutpat. Cras mollis, ex quis tristique ullamcorper, risus est ultrices nisi, eu cursus mauris turpis maximus metus. Ut eu luctus dolor. Donec consequat molestie justo in efficitur. Phasellus eleifend at dui sit amet varius. Aenean ornare elit a nunc dapibus, vel fermentum nibh rutrum. Nulla imperdiet auctor est ac hendrerit. Vestibulum feugiat eleifend enim non tincidunt. Morbi ut auctor libero. Aliquam quis vulputate velit, a sodales lacus. Curabitur ante augue, tincidunt nec vestibulum a, maximus in risus. Fusce ultrices, est sit amet ornare ullamcorper, sapien lorem tincidunt velit, at ultricies diam risus vel justo. Nulla eros est, efficitur quis est ut, sollicitudin scelerisque purus. Nullam sed lacus nisl. Fusce scelerisque augue aliquet quam porta, ut euismod purus malesuada."
  
  buffer := new(bytes.Buffer)
  reader := NewReaderFramer(buffer)
  writer := NewWriterFramer(buffer)
  
  if err := writer.Write([]byte(message)); err != nil {
    t.Errorf("Could not write message: %v", err)
  }
  if err := writer.Write([]byte(message)); err != nil {
    t.Errorf("Could not write message: %v", err)
  }
  if err := writer.Write([]byte(message)); err != nil {
    t.Errorf("Could not write message: %v", err)
  }
  
  for buffer.Len() > 0 {
    fmt.Printf("--- %v ---\n", buffer.Len())
    if m, err := reader.Read(); err != nil {
      t.Errorf("Could not read message: %v", err)
    }else if len(m) < 1 {
      t.Errorf("Invalid number of messages read: %d < 1", len(m))
    }else if string(m[0]) != message {
      t.Errorf("Invalid message data: %+v != %+v", string(m[0]), message)
    }else{
      for _, e := range m {
        fmt.Printf("Received: %q\n", string(e))
      }
    }
  }
  
}

