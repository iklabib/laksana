using System.Text.Json.Serialization;

public class BuildResult 
{
    [JsonPropertyName("status")]
    public StatusType Status {get; set;}

    [JsonPropertyName("message")]
    public string Message {get; set;}

    [JsonPropertyName("compilation_errors")]
    public CompileError[] CompilatioErrors {get; set;}
}

public class CompileError
{
    [JsonPropertyName("filename")]
    public string Filename {get; set;}

    [JsonPropertyName("message")]
    public string Message {get; set;}

    [JsonPropertyName("line")]
    public int Line {get; set;}

    [JsonPropertyName("character")]
    public int Character {get; set;}
}

public class Submission 
{
    [JsonPropertyName("filename")] 
    public string Filename {get; set;}

    [JsonPropertyName("path")] 
    public string Path {get; set;}

    [JsonPropertyName("src")] 
    public string SourceCode {get; set;}
}

public class TestResult
{
    [JsonPropertyName("status")] 
    public string Status {get; set;}
    public string Name {get; set;}
    public string Output {get; set;}
    public int Order {get; set;}
}

public enum StatusType 
{
    OK = 0,
    ERROR,
    INTERNAL_ERROR,
}