using System.Text.Json;

string command =  args[0];
string assemblyName = args[1];
var testMan = new TestManager();
switch (command)
{
    case "build":
        string? lines = Console.ReadLine();
        var buildResult = testMan.Build(assemblyName, lines);
        string serialized = JsonSerializer.Serialize(buildResult);
        Console.WriteLine(serialized);
    break;

    case "execute":
        var testResult = testMan.Run(assemblyName);
        Console.WriteLine(testResult);
    break;
}