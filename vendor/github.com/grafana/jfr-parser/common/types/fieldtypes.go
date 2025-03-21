package types

type FieldClass string

const (
	Boolean              FieldClass = "boolean"
	Byte                 FieldClass = "byte"
	Char                 FieldClass = "char"
	Short                FieldClass = "short"
	Int                  FieldClass = "int"
	Long                 FieldClass = "long"
	Float                FieldClass = "float"
	Double               FieldClass = "double"
	String               FieldClass = "java.lang.String"
	Class                FieldClass = "java.lang.Class"
	Thread               FieldClass = "java.lang.Thread"
	ClassLoader          FieldClass = "jdk.types.ClassLoader"
	CodeBlobType         FieldClass = "jdk.types.CodeBlobType"
	FlagValueOrigin      FieldClass = "jdk.types.FlagValueOrigin"
	FrameType            FieldClass = "jdk.types.FrameType"
	G1YCType             FieldClass = "jdk.types.G1YCType"
	GCName               FieldClass = "jdk.types.GCName"
	Method               FieldClass = "jdk.types.Method"
	Module               FieldClass = "jdk.types.Module"
	NarrowOopMode        FieldClass = "jdk.types.NarrowOopMode"
	NetworkInterfaceName FieldClass = "jdk.types.NetworkInterfaceName"
	Package              FieldClass = "jdk.types.Package"
	StackFrame           FieldClass = "jdk.types.StackFrame"
	StackTrace           FieldClass = "jdk.types.StackTrace"
	Symbol               FieldClass = "jdk.types.Symbol"
	ThreadState          FieldClass = "jdk.types.ThreadState"
	InflateCause         FieldClass = "jdk.types.InflateCause"
	GCCause              FieldClass = "jdk.types.GCCause"
	CompilerPhaseType    FieldClass = "jdk.types.CompilerPhaseType"
	ThreadGroup          FieldClass = "jdk.types.ThreadGroup"
	GCThresholdUpdater   FieldClass = "jdk.types.GCThresholdUpdater"
	MetaspaceObjectType  FieldClass = "jdk.types.MetaspaceObjectType"
	ExecutionMode        FieldClass = "datadog.types.ExecutionMode"
	VMOperationType      FieldClass = "jdk.types.VMOperationType"
	G1HeapRegionType     FieldClass = "jdk.types.G1HeapRegionType"
	GCWhen               FieldClass = "jdk.types.GCWhen"
	ReferenceType        FieldClass = "jdk.types.ReferenceType"
	MetadataType         FieldClass = "jdk.types.MetadataType"
	LogLevel             FieldClass = "profiler.types.LogLevel"
	AttributeValue       FieldClass = "profiler.types.AttributeValue"
)
