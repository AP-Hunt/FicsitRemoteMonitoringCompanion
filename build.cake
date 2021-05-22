
var target = Argument("target", "Test");
var configuration = Argument("configuration", "Release");

var prometheusVersion = "2.27.1";
var grafanaVersion = "7.5.7";

//////////////////////////////////////////////////////////////////////
// TASKS
//////////////////////////////////////////////////////////////////////

Task("Clean")
    .WithCriteria(c => HasArgument("rebuild"))
    .Does(() =>
{
    CleanDirectory($"./bin/{configuration}");
});

Task("Build")
    .IsDependentOn("Clean")
    .IsDependentOn("Prometheus")
    .IsDependentOn("Grafana")
    .Does(() =>
{
    DotNetCoreBuild("./FicsitRemoteMonitoringCompanion.sln", new DotNetCoreBuildSettings
    {
        Configuration = configuration,
        OutputDirectory = $"./bin/{configuration}"
    });

    CopyFile($"./Externals/Prometheus/prometheus-{prometheusVersion}.windows-amd64/prometheus.exe", $"./bin/{configuration}/prometheus.exe");
    CopyDirectory($"./Externals/Grafana/grafana-{grafanaVersion}", $"./bin/{configuration}/grafana");
});

Task("Test")
    .IsDependentOn("Build")
    .Does(() =>
{
    DotNetCoreTest("./FicsitRemoteMonitoringCompanion.sln", new DotNetCoreTestSettings
    {
        Configuration = configuration,
        NoBuild = true,
    });
});

Task("Prometheus")
    .IsDependentOn("ExternalsDir")
    .Does(() => 
{
    System.IO.Directory.CreateDirectory("./Externals/Prometheus");

    if(!System.IO.Directory.Exists($"./Externals/Prometheus/prometheus-{prometheusVersion}.windows-amd64"))
    {
        Information($"Downloading Prometheus {prometheusVersion}");
        DownloadFile(
          $"https://github.com/prometheus/prometheus/releases/download/v{prometheusVersion}/prometheus-{prometheusVersion}.windows-amd64.zip",
          "./Externals/Prometheus/prometheus.zip"
        );
        Unzip("./Externals/Prometheus/prometheus.zip", "./Externals/Prometheus");
    } else 
    {
        Information("Skipping Prometheus download because it has already been downloaded");    
    }
});

Task("Grafana")
    .IsDependentOn("ExternalsDir")
    .Does(() => 
{
    System.IO.Directory.CreateDirectory("./Externals/Grafana");

    if(!System.IO.Directory.Exists($"./Externals/Grafana/grafana-{grafanaVersion}"))
    {
        Information($"Downloading Grafana {grafanaVersion}");
        DownloadFile(
          $"https://dl.grafana.com/oss/release/grafana-{grafanaVersion}.windows-amd64.zip ",
          "./Externals/Grafana/grafana.zip"
        );
        Unzip("./Externals/Grafana/grafana.zip", "./Externals/Grafana");
    } else 
    {
        Information("Skipping Grafana download because it has already been downloaded");    
    }
});

Task("ExternalsDir")
    .Does(() => 
{
    System.IO.Directory.CreateDirectory("./Externals/");
});

//////////////////////////////////////////////////////////////////////
// EXECUTION
//////////////////////////////////////////////////////////////////////

RunTarget(target);