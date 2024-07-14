# Laksana
An experimental code executor, this project is designed to run within a container and in conjunction with other applications, rather than as a standalone application. Currently, only Docker is supported. Other container engines, such as Podman, might work in rootful mode.

On AppArmor enforced systems (e.g. Ubuntu), Laksana requires the installation of an AppArmor profile. If AppArmor is not enforced, skip the first step.
```bash
$ ./chore.bash apparmor
$ ./chore.bash build
# apparmor option is required on AppArmor enforced systems
$ ./chore.bash run --apparmor
```
Laksana should be serving by now. Give it a test.
```bash
$ curl --request POST \
  --url http://localhost:31415/run \
  --header 'content-type: application/json' \
  --data '{
  "type": "csharp",
  "src": [
    {
      "filename": "Program.cs",
      "src": "using System;\n\t\t\t\t\t\npublic class Program\n{\n\tpublic static void Main(string[] args)\n\t{\n\t\tforeach (string arg in args) {\n\t\t\tConsole.WriteLine(arg);\n\t\t}\n\t}\n}"
    }
  ],
  "src_test": "using System.Reflection;\nusing NUnit.Framework;\nusing NUnitLite;\n\nnamespace test;\n\n[TestFixture]\npublic class Tests\n{\n    public static int Main(string[] args) \n    {\n        return new AutoRun(Assembly.GetExecutingAssembly()).Execute([\"/test:Tests.Test1\"]);\n    }\n\n    [SetUp]\n    public void Setup()\n    {\n    }\n\n    [Test, Order(1)]\n    public void Test1()\n    {\n        Assert.Equals(1, 2);\n    }\n\n    [Test, Order(2)]\n    public void Test2()\n    {\n        Assert.Pass(\"OKEY\");\n    }\n}"
}'
```
Expected response would be as follows.
```json
{
  "success": true,
  "message": "",
  "builds": null,
  "tests": [
    {
      "status": "FAILED",
      "name": "Test1",
      "output": "",
      "order": 1
    },
    {
      "status": "PASSED",
      "name": "Test2",
      "output": "OKEY",
      "order": 2
    }
  ]
}
```
