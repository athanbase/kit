### go中protocol buffer的基本使用(proto3)

这个教程介绍了go中protocol buffers基于proto3版本的基础使用方法。通过创建一个简单的示例应用来介绍：

- 在.proto文件中定义message格式

- 使用 protocol buffer编译

- 使用 Go protocol buffer API读写messages

这不是一个全面的go使用protocol buffer手册。更多详细信息可以查阅[Protocol Buffer Language Guide](https://developers.google.com/protocol-buffers/docs/proto3), [Go API Reference](https://pkg.go.dev/google.golang.org/protobuf/proto), [Go Generated Code Guide](https://developers.google.com/protocol-buffers/docs/reference/go-generated), 以及[Encodeing Reference](https://developers.google.com/protocol-buffers/docs/encoding)。

### 选择protocol buffers的理由

我们使用一个很简单的"address book"应用作为示例：可以从一个文件中读取联系人的详细联系方式。每条联系人记录包含：name、ID、email address, contact phone number。

你是如何序列化/检索这样的结构化数据? 有一下几种方式解决这个问题: 

- 使用gobs序列化Go结构体。在单纯Go开发环境中是个很好的方式，如果要跟其它语言交换信息这种方式就行不通了。

- 你可以发明一种将数据编码为单个字符串的临时方法：例如编码成"12:3:-23:67”这样的4个整数。虽然需要编写一次性编码和解码代码，但是这种方法简单、灵活，解码成本很小。对于结构简单的数据很有效。

- 把数据序列化成XML格式。XML是人类（一种）可读的，多种语言都有现成的解析库，这种方法可能非常有吸引力。 如果要与其他应用程序/项目共享数据，这可能是一个不错的选择。但是XML占用大量空间，对它进行编码/解码会给应用程序带来巨大的性能损失。 遍历XML DOM树比遍历类中的简单字段要复杂得多。

对于这些问题Protocol buffers是一个灵活、高效、自动化的解决方案。通过protocol buffers，你只需写一个.proto文件描述要存储数据的结构。之后protocol buffer编译器会创建一个把protocol buffer数据自动编码和解析成二进制格式的对象。生成的对象提供了字段的读取和修改方法。Importantly, the protocol buffer format supports the idea of extending the format over time in such a way that the code can still read data encoded with the old format.


[示例代码地址](https://github.com/protocolbuffers/protobuf/tree/master/examples)


### 定义协议格式

以一个.proto为开始创建你的address book应用。.proto文件中的定义很简单：为每个要序列化的结构体添加message,然后为每个message的字段设置名字和类型。在这个例子中包含message定义的.proto文件名叫addressbook.proto。

.proto文件以包定义开头防止不同项目之间的名字冲突

go_package选项定义了包含本文件生成所有代码的导入路径。go报名是导入路径中最后一段。本示例使用的报名是: tutorialpb。

接下来是定义message。一个message只是一组类型字段的汇总。许多数据基础类型可以用与字段类型包括: bool, int32, float, double, string。 也可以使用自定义的message作为字段类型。


```proto
syntax = "proto3";
package tutorial;

import "google/protobuf/timestamp.proto";

option go_package = "github.com/protocolbuffers/protobuf/examples/go/tutorialpb";


message Person {
  string name = 1;
  int32 id = 2;  // Unique ID number for this person.
  string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
  }

  repeated PhoneNumber phones = 4;

  google.protobuf.Timestamp last_updated = 5;
}

// Our address book file is just one of these.
message AddressBook {
  repeated Person people = 1;
}

```

In the above example, the Person message contains PhoneNumber messages, while the AddressBook message contains Person messages. You can even define message types nested inside other messages – as you can see, the PhoneNumber type is defined inside Person. You can also define enum types if you want one of your fields to have one of a predefined list of values – here you want to specify that a phone number can be one of MOBILE, HOME, or WORK.

The " = 1", " = 2" markers on each element identify the unique "tag" that field uses in the binary encoding. Tag numbers 1-15 require one less byte to encode than higher numbers, so as an optimization you can decide to use those tags for the commonly used or repeated elements, leaving tags 16 and higher for less-commonly used optional elements. Each element in a repeated field requires re-encoding the tag number, so repeated fields are particularly good candidates for this optimization.

If a field value isn't set, a default value is used: zero for numeric types, the empty string for strings, false for bools. For embedded messages, the default value is always the "default instance" or "prototype" of the message, which has none of its fields set. Calling the accessor to get the value of a field which has not been explicitly set always returns that field's default value.

If a field is repeated, the field may be repeated any number of times (including zero). The order of the repeated values will be preserved in the protocol buffer. Think of repeated fields as dynamically sized arrays.

You'll find a complete guide to writing .proto files – including all the possible field types – in the [Protocol Buffer Language Guide](https://developers.google.com/protocol-buffers/docs/proto3). Don't go looking for facilities similar to class inheritance, though – protocol buffers don't do that.

### 编译 protocol buffers

Now that you have a .proto, the next thing you need to do is generate the classes you'll need to read and write AddressBook (and hence Person and PhoneNumber) messages. To do this, you need to run the protocol buffer compiler protoc on your .proto:

1. If you haven't installed the compiler, [download the package](https://developers.google.com/protocol-buffers/docs/downloads) and follow the instructions in the README.

2. Run the following command to install the Go protocol buffers plugin:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

The compiler plugin protoc-gen-go will be installed in $GOBIN, defaulting to $GOPATH/bin. It must be in your $PATH for the protocol compiler protoc to find it.

3. Now run the compiler, specifying the source directory (where your application's source code lives – the current directory is used if you don't provide a value), the destination directory (where you want the generated code to go; often the same as $SRC_DIR), and the path to your .proto. In this case, you would invoke:

```sh
protoc -I=`pwd` --go_out=`pwd` `pwd`/addressbook.proto
```