using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.Text;

try 
{
    string? input = Console.ReadLine();
    if (string.IsNullOrEmpty(input)) 
    {
        Console.WriteLine("[]");
        return;
    }

    var submissions = JsonSerializer.Deserialize<List<Submission>>(input);
    if (submissions == null || submissions.Count() == 0)
    {
        Console.WriteLine("[]");
        return;
    }

    // implicit using workaround
    // https://github.com/dotnet/roslyn/issues/58119
    var implicingUsings = CSharpSyntaxTree.ParseText(SourceText.From("""
    global using global::System;
    global using global::System.IO;
    global using global::System.Linq;
    global using global::System.Net.Http;
    global using global::System.Threading;
    global using global::System.Threading.Tasks;
    global using global::System.Collections.Generic;
    """))
    .WithFilePath("GlobalUsings.cs");

    var syntaxTree = submissions.Select(el => { 
            var source = SourceText.From(el.Content);
            return CSharpSyntaxTree.ParseText(source).WithFilePath(el.Filename);
        })
        .Prepend(implicingUsings);

    var references = ((string)AppContext.GetData("TRUSTED_PLATFORM_ASSEMBLIES")!)
                     .Split(Path.PathSeparator)
                     .Where(el => !string.IsNullOrEmpty(el))
                     .Select(el => MetadataReference.CreateFromFile(el));

    CSharpCompilationOptions compilationOptions = new CSharpCompilationOptions(OutputKind.ConsoleApplication)
                                                    .WithOptimizationLevel(OptimizationLevel.Release)
                                                    .WithPlatform(Platform.X64)
                                                    .WithWarningLevel(1);

    var compilation = CSharpCompilation.Create("LittleRosie")
                                    .WithOptions(compilationOptions)
                                    .AddReferences(references)
                                    .AddSyntaxTrees(syntaxTree);


    List<CompileError> compilationError = new();
    var errors = compilation.GetDiagnostics().Where(el => el.Severity == DiagnosticSeverity.Error);
    if (errors.Count() == 0)
    {
        Console.WriteLine("[]");
        return;
    }

    foreach(var d in errors)
    {
        var location = d.Location;
        var line = location.GetLineSpan().StartLinePosition;
        var compileError = new CompileError {
            Filename = location.SourceTree?.FilePath ?? "",
            Message = d.GetMessage(),
            Line = line.Line + 1,
            Column = line.Character,
        };
        compilationError.Add(compileError);
    }

    Console.WriteLine(JsonSerializer.Serialize(compilationError));
} 
catch(Exception) 
{
    Console.WriteLine("[]");
}

public class CompileError
{
    public string Filename {get; set;} = "";
    public string Message {get; set;} = "";
    public int Line {get; set;}
    public int Column {get; set;}
}

public class Submission 
{
    [JsonPropertyName("filename")] public string Filename {get; set;} = "";
    [JsonPropertyName("content")] public string Content {get; set;} = "";
}