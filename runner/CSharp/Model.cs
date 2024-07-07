using System.Text.Json.Serialization;
using System.Xml;

public class BuildResult 
{
    [JsonPropertyName("status")]
    public StatusType Status {get; set;}

    [JsonPropertyName("message")]
    public string Message {get; set;}

    [JsonPropertyName("compilation_erros")]
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

    [JsonPropertyName("content")] 
    public string Content {get; set;}

    [JsonPropertyName("is_test")] 
    public bool IsTest {get; set;} = false;
}

public class TestResult
{
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