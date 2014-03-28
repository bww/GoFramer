// 
// Copyright (c) 2014 Brian W. Wolter, All rights reserved.
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
// Developed in New York City
// 

package framer

import (
  "fmt"
  "bytes"
  "testing"
)

func TestReadWrite(t *testing.T) {
  message := "Hello, this is the message"
  
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
  
  fmt.Println(buffer.Bytes())
  
  for buffer.Len() > 0 {
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

