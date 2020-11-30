
package zrun


var embeddedPythonProlog string = string ([]byte {
  0x23, 0x21, 0x2f, 0x64, 0x65, 0x76, 0x2f, 0x6e, 0x75, 0x6c, 0x6c, 0x0a, 0x0a, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23,
  0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x0a, 0x0a, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x28, 0x5a, 0x2c, 0x20, 0x5f, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x20, 0x3d, 0x20, 0x5f, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x5a, 0x2e, 0x6f, 0x73, 0x20, 0x3d, 0x20, 0x5f, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x20, 0x28, 0x22, 0x6f, 0x73, 0x22, 0x29, 0x0a, 0x09, 0x5a, 0x2e, 0x73, 0x79, 0x73, 0x20, 0x3d, 0x20, 0x5f, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x20, 0x28, 0x22, 0x73, 0x79, 0x73, 0x22, 0x29, 0x0a, 0x09, 0x5a, 0x2e, 0x73, 0x68, 0x75, 0x74, 0x69, 0x6c, 0x20, 0x3d, 0x20, 0x5f,
  0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x20, 0x28, 0x22, 0x73, 0x68, 0x75, 0x74, 0x69, 0x6c, 0x22, 0x29, 0x0a, 0x09, 0x5a, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x3d, 0x20, 0x5f, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x5f, 0x20, 0x28, 0x22, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x28, 0x5f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x3d, 0x20, 0x5f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x5f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x5f, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x77, 0x69, 0x74, 0x68, 0x20, 0x28, 0x22, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x22, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x3d, 0x20, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5b, 0x38, 0x3a, 0x5d, 0x0a, 0x09, 0x09,
  0x65, 0x6c, 0x73, 0x65, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x61, 0x73, 0x73, 0x65, 0x72, 0x74, 0x20, 0x46, 0x61, 0x6c, 0x73, 0x65, 0x2c, 0x20, 0x28, 0x22, 0x5b, 0x38, 0x33, 0x63, 0x65, 0x63, 0x38, 0x34, 0x39, 0x5d, 0x20, 0x20, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x20, 0x60, 0x25, 0x73, 0x60, 0x22, 0x20, 0x25, 0x20, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x29, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x5f, 0x64, 0x69, 0x63, 0x74, 0x5f, 0x5f, 0x5b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5d, 0x20, 0x3d, 0x20, 0x5f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x0a, 0x09, 0x0a, 0x09,
  0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x7a, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c,
  0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x5f, 0x7a, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5a, 0x2e, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
  0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x7a, 0x65, 0x78, 0x65, 0x63, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x5f, 0x7a, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e,
  0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5a, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x5f, 0x7a, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65,
  0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x77, 0x69, 0x74, 0x68, 0x20, 0x28, 0x22, 0x3a, 0x3a, 0x22, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5a, 0x2e, 0x70, 0x61, 0x6e, 0x69, 0x63, 0x20, 0x28, 0x30, 0x78, 0x62, 0x64, 0x31, 0x36, 0x34, 0x31, 0x63, 0x37, 0x2c, 0x20, 0x22, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x3a, 0x20, 0x60, 0x25, 0x73, 0x60, 0x22, 0x2c, 0x20, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x29, 0x0a, 0x09, 0x09, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65,
  0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x20, 0x3d, 0x20, 0x5b, 0x22, 0x5b, 0x7a, 0x2d, 0x72, 0x75, 0x6e, 0x5d, 0x22, 0x2c, 0x20, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x5d, 0x0a, 0x09, 0x09, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x20, 0x28, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x7b, 0x7d, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x28, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x29,
  0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x46, 0x61, 0x6c, 0x73, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70,
  0x74, 0x6c, 0x65, 0x74, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5a, 0x2e, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x20, 0x28, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
  0x6e, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5a, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2c, 0x20, 0x2a, 0x2a, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a,
  0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x20, 0x28, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x20, 0x3d, 0x20, 0x5b, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5d, 0x0a, 0x09, 0x09, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x20, 0x28, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
  0x65, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x7b, 0x7d, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x28, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x54, 0x72, 0x75, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2c, 0x20, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x20, 0x3d, 0x20, 0x54, 0x72, 0x75, 0x65, 0x2c, 0x20, 0x5f, 0x70, 0x61, 0x6e,
  0x69, 0x63, 0x20, 0x3d, 0x20, 0x54, 0x72, 0x75, 0x65, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x76, 0x70, 0x65, 0x0a, 0x09, 0x09, 0x65, 0x6c, 0x73, 0x65, 0x20, 0x3a,
  0x0a, 0x09, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x76, 0x65, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x20, 0x3d, 0x20, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x28, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x50, 0x5f, 0x57, 0x41, 0x49, 0x54, 0x2c, 0x20, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x29, 0x0a, 0x09, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x70, 0x61,
  0x6e, 0x69, 0x63, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x20, 0x21, 0x3d, 0x20, 0x30, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x09, 0x5a, 0x2e, 0x70, 0x61, 0x6e, 0x69, 0x63, 0x20, 0x28, 0x30, 0x78, 0x33, 0x63, 0x31, 0x34, 0x62, 0x39, 0x61, 0x30, 0x2c, 0x20, 0x22, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x20, 0x60, 0x25, 0x73, 0x60, 0x20, 0x60, 0x25, 0x73, 0x60, 0x20, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x3a, 0x20, 0x25, 0x64, 0x22, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5b, 0x30, 0x5d, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5b, 0x31, 0x3a, 0x5d, 0x2c, 0x20, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f,
  0x6d, 0x65, 0x29, 0x0a, 0x09, 0x09, 0x65, 0x6c, 0x73, 0x65, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x20, 0x3d, 0x20, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x28, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x50, 0x5f, 0x4e, 0x4f, 0x57, 0x41, 0x49, 0x54, 0x2c, 0x20, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x29, 0x0a, 0x09, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x70, 0x61, 0x6e, 0x69, 0x63, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x20, 0x3c, 0x3d, 0x20, 0x30, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09,
  0x09, 0x5a, 0x2e, 0x70, 0x61, 0x6e, 0x69, 0x63, 0x20, 0x28, 0x30, 0x78, 0x33, 0x36, 0x37, 0x33, 0x37, 0x64, 0x34, 0x38, 0x2c, 0x20, 0x22, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x20, 0x60, 0x25, 0x73, 0x60, 0x20, 0x60, 0x25, 0x73, 0x60, 0x20, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x28, 0x25, 0x64, 0x29, 0x3a, 0x20, 0x25, 0x73, 0x22, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5b, 0x30, 0x5d, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5b, 0x31, 0x3a, 0x5d, 0x2c, 0x20, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x2c, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x28, 0x5f, 0x6f, 0x75, 0x74,
  0x63, 0x6f, 0x6d, 0x65, 0x29, 0x29, 0x0a, 0x09, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5f, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x5f, 0x30, 0x20, 0x28, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x5f, 0x64, 0x65, 0x73,
  0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x5f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x76, 0x70, 0x65, 0x0a, 0x09, 0x09, 0x65, 0x6c, 0x73, 0x65, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x76, 0x65, 0x0a, 0x09, 0x09, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x20, 0x28, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2c, 0x20, 0x5f, 0x65, 0x6e,
  0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x65,
  0x72, 0x72, 0x6f, 0x72, 0x20, 0x28, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x27, 0x65, 0x65, 0x27, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x20, 0x28, 0x5f,
  0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3a, 0x20, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x27, 0x77, 0x77, 0x27, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69,
  0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x20, 0x28, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x6f, 0x72, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3a, 0x20, 0x72, 0x65, 0x74, 0x75,
  0x72, 0x6e, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x27, 0x69, 0x69, 0x27, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x20, 0x28, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x69, 0x66, 0x20, 0x6e,
  0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x6f, 0x72, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x6f, 0x72, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3a, 0x20, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x27, 0x64, 0x64, 0x27, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20,
  0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x5f, 0x73, 0x6c, 0x75, 0x67, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x28, 0x22, 0x5b, 0x7a, 0x2d, 0x72, 0x75, 0x6e, 0x3a, 0x25, 0x30, 0x38, 0x64, 0x5d, 0x20, 0x5b, 0x25, 0x73, 0x5d, 0x20, 0x5b, 0x25, 0x30,
  0x38, 0x78, 0x5d, 0x20, 0x20, 0x22, 0x20, 0x25, 0x20, 0x28, 0x5a, 0x2e, 0x70, 0x69, 0x64, 0x2c, 0x20, 0x5f, 0x73, 0x6c, 0x75, 0x67, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x29, 0x29, 0x20, 0x2b, 0x20, 0x28, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x20, 0x25, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x2b, 0x20, 0x22, 0x5c, 0x6e, 0x22, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x20, 0x28, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x65, 0x78, 0x69, 0x74, 0x20, 0x28, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x70, 0x61, 0x6e, 0x69, 0x63, 0x20, 0x28, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c,
  0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x2a, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x20, 0x28, 0x27, 0x21, 0x21, 0x27, 0x2c, 0x20, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x2c, 0x20, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x29, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x65, 0x78, 0x69, 0x74, 0x20, 0x28, 0x31, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x40, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x73, 0x6c, 0x65, 0x65, 0x70, 0x20, 0x28, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72,
  0x76, 0x61, 0x6c, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x5a, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x73, 0x6c, 0x65, 0x65, 0x70, 0x20, 0x28, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x29, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x5a, 0x2e, 0x70, 0x69,
  0x64, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x67, 0x65, 0x74, 0x70, 0x69, 0x64, 0x20, 0x28, 0x29, 0x0a, 0x09, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x6f, 0x73, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x0a, 0x09, 0x5a, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5b, 0x22, 0x5a, 0x52, 0x55, 0x4e, 0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x22, 0x5d, 0x0a, 0x09, 0x5a, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5b, 0x22,
  0x5a, 0x52, 0x55, 0x4e, 0x5f, 0x57, 0x4f, 0x52, 0x4b, 0x53, 0x50, 0x41, 0x43, 0x45, 0x22, 0x5d, 0x0a, 0x09, 0x5a, 0x2e, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5b, 0x22, 0x5a, 0x52, 0x55, 0x4e, 0x5f, 0x46, 0x49, 0x4e, 0x47, 0x45, 0x52, 0x50, 0x52, 0x49, 0x4e, 0x54, 0x22, 0x5d, 0x0a, 0x09, 0x0a, 0x09, 0x5a, 0x2e, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x0a, 0x09, 0x5a, 0x2e, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x20, 0x3d, 0x20, 0x5a, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x0a, 0x09, 0x5a, 0x2e, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x20,
  0x3d, 0x20, 0x5a, 0x2e, 0x73, 0x79, 0x73, 0x2e, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x0a, 0x09, 0x0a, 0x09, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3d, 0x20, 0x54, 0x72, 0x75, 0x65, 0x0a, 0x09, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3d, 0x20, 0x54, 0x72, 0x75, 0x65, 0x0a, 0x09, 0x5a, 0x2e, 0x6c, 0x6f, 0x67, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x3d, 0x20, 0x46, 0x61, 0x6c, 0x73, 0x65, 0x0a, 0x09, 0x0a, 0x09, 0x23, 0x23, 0x20, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d,
  0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x0a, 0x09, 0x0a, 0x09, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x20, 0x5a, 0x0a, 0x0a, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23,
  0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x0a, 0x0a, 0x69, 0x66, 0x20, 0x5f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x5f, 0x20, 0x3d, 0x3d,
  0x20, 0x22, 0x5f, 0x5f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x5f, 0x22, 0x20, 0x3a, 0x0a, 0x09, 0x0a, 0x23, 0x21, 0x09, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x73, 0x79, 0x73, 0x0a, 0x23, 0x21, 0x09, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x6f, 0x73, 0x0a, 0x09, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x20, 0x28, 0x29, 0x20, 0x3a, 0x0a, 0x09, 0x09, 0x61, 0x73, 0x73, 0x65, 0x72, 0x74, 0x20, 0x46, 0x61, 0x6c, 0x73, 0x65, 0x2c, 0x20, 0x22, 0x5b, 0x63, 0x39, 0x32, 0x62, 0x62, 0x35, 0x38, 0x35, 0x5d, 0x22, 0x0a, 0x09, 0x0a, 0x09, 0x7a, 0x72, 0x75, 0x6e, 0x20, 0x3d, 0x20, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x5f, 0x5f, 0x69, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x28, 0x5f, 0x5f, 0x7a, 0x72, 0x75, 0x6e, 0x29, 0x0a, 0x09, 0x0a, 0x65,
  0x6c, 0x73, 0x65, 0x20, 0x3a, 0x0a, 0x09, 0x0a, 0x09, 0x61, 0x73, 0x73, 0x65, 0x72, 0x74, 0x20, 0x46, 0x61, 0x6c, 0x73, 0x65, 0x2c, 0x20, 0x28, 0x22, 0x5b, 0x62, 0x31, 0x31, 0x34, 0x35, 0x38, 0x61, 0x37, 0x5d, 0x20, 0x20, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x20, 0x60, 0x25, 0x73, 0x60, 0x22, 0x20, 0x25, 0x20, 0x5f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x5f, 0x29, 0x0a, 0x0a, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23,
  0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x0a, 0x0a,
})

