using System.Reflection;
using NUnit.Framework;
using NUnitLite;

namespace test;

[TestFixture]
public class Tests
{
    public static int Main(string[] args) 
    {
        return new AutoRun(Assembly.GetExecutingAssembly()).Execute(["/test:Tests.Test1"]);
    }

    [SetUp]
    public void Setup()
    {
    }

    [Test, Order(1)]
    public void Test1()
    {
        Assert.Equals(1, 2);
    }

    [Test, Order(2)]
    public void Test2()
    {
        Assert.Pass("OKE BRUH");
    }
}