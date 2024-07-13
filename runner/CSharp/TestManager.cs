using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.Text;
using NUnit;
using NUnit.Engine;
using System.Reflection;
using System.Text.Json;
using System.Xml;

public class TestManager 
{ 
    private string assemblyName = "main.dll";

    public string Run() 
    {
        using var engine = TestEngineActivator.CreateInstance();
        var assembly = Assembly.LoadFrom(assemblyName);
        var package = new TestPackage(assemblyName);

        using var runner = engine.GetRunner(package);
        XmlNode testResult = runner.Run(null, TestFilter.Empty);
        return parseTestResult(testResult);
    }

    public BuildResult Build(string lines) 
    {
        try 
        {
            if (string.IsNullOrEmpty(lines)) 
            { 
                return new BuildResult 
                {
                    Status = StatusType.INTERNAL_ERROR,
                    Message = "empty input",
                };
            }

            var submissions = JsonSerializer.Deserialize<List<Submission>>(lines) ?? new List<Submission>();
            if (submissions.Count() == 0)
            {
                return new BuildResult 
                {
                    Status = StatusType.INTERNAL_ERROR,
                    Message = "serialization failed",
                };
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

            List<SyntaxTree> syntaxTree = submissions.Where(el => !string.IsNullOrEmpty(el.SourceCode))
                            .Select(el => { 
                                var source = SourceText.From(el.SourceCode);
                                return CSharpSyntaxTree.ParseText(source).WithFilePath(el.Filename);
                            }).Prepend(implicingUsings).ToList();

            var references = ((string)AppContext.GetData("TRUSTED_PLATFORM_ASSEMBLIES")!)
                            .Split(Path.PathSeparator)
                            .Where(el => !string.IsNullOrEmpty(el))
                            .Select(el => MetadataReference.CreateFromFile(el));

            if (references.Count() == 0) 
            {
                return new BuildResult 
                {
                    Status = StatusType.INTERNAL_ERROR,
                    Message = "no trusted assemblies",
                };
            }

            var compilationOptions = new CSharpCompilationOptions(OutputKind.ConsoleApplication)
                                                            .WithOptimizationLevel(OptimizationLevel.Release)
                                                            .WithPlatform(Platform.X64)
                                                            .WithWarningLevel(1);

            var compilation = CSharpCompilation.Create(assemblyName)
                                            .WithOptions(compilationOptions)
                                            .AddReferences(references)
                                            .AddSyntaxTrees(syntaxTree);

            var emitResult = compilation.Emit(assemblyName);
            var errors = emitResult.Diagnostics.Where(el => el.Severity == DiagnosticSeverity.Error);
            if (emitResult.Success && errors.Count() == 0)
            {
                return new BuildResult { Status = StatusType.OK };
            }
            else
            {
                return new BuildResult 
                { 
                    Status = StatusType.ERROR, 
                    CompilatioErrors = Diagnostics(errors),
                };
            }
        } 
        catch(Exception e) 
        {
            return new BuildResult 
            { 
                Status = StatusType.INTERNAL_ERROR,
                Message = e.StackTrace ?? e.Message,
            };
        }
    }

    public CompileError[] Diagnostics(IEnumerable<Diagnostic> errors)
    {
        var compilationError = new List<CompileError>();
        foreach(var d in errors)
        {
            var location = d.Location;
            var line = location.GetLineSpan().StartLinePosition;
            var compileError = new CompileError {
                Filename = location.SourceTree?.FilePath ?? "",
                Message = d.GetMessage(),
                Line = line.Line,
                Character = line.Character,
            };
            compilationError.Add(compileError);
        }

        return compilationError.ToArray();
    }

    private string parseTestResult(XmlNode node)
    {
        var testCaseNodes = node.SelectNodes("//test-case");
        if (testCaseNodes == null)
        {
            return "[]";
        }

        var testResults = new List<TestResult>();
        foreach (XmlNode testCase in testCaseNodes)
        {
            string status = testCase.GetAttribute("result").ToUpper();

            int order = 0;
            var orderNode = testCase.SelectSingleNode("properties/property[@name='Order']");
            if (orderNode != null)
            {
                order = int.Parse(orderNode.GetAttribute("value"));
            }

            string message;
            if (status == "PASSED")
            {
                message = testCase?.SelectSingleNode("reason/message")?.InnerText ?? "";
            }
            else
            {
                message = testCase?.SelectSingleNode("reason/failure")?.InnerText ?? "";
            }

            var result = new TestResult
            {
                Status = testCase.GetAttribute("result").ToUpper(),
                Name = testCase.GetAttribute("name"),
                Output = message,
                Order = order,
            };

            testResults.Add(result);
        }

        return JsonSerializer.Serialize(testResults);
    }
}
