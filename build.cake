#addin nuget:?package=semver&version=2.0.4
#addin nuget:?package=Cake.Git&version=1.0.1
#addin nuget:?package=Cake.Npm&version=1.0.0

using Semver; 
using Cake.Core.IO;

var target = Argument("target", "Test");
var configuration = Argument("configuration", "Release");

var prometheusVersion = "2.27.1";
var grafanaVersion = "7.5.7";
var gitchglogVersion = "0.14.2";

var msBuildSettings = new DotNetCoreMSBuildSettings();

if(configuration == "Release")
{
    string version = System.IO.File.ReadAllText("version.txt");
    var semVersion = SemVersion.Parse(version, false);
    
    msBuildSettings.Properties["AssemblyVersionNumber"] = new string[]{semVersion.ToString()};
}

//////////////////////////////////////////////////////////////////////
// TASKS
//////////////////////////////////////////////////////////////////////

Task("Clean")
    .WithCriteria(c => HasArgument("rebuild"))
    .Does(() =>
{
    CleanDirectory($"./bin/{configuration}");
    CleanDirectory($"./bin/publish/{configuration}");
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
        OutputDirectory = $"./bin/{configuration}",
        MSBuildSettings = msBuildSettings
    });

    if(configuration == "Release")
    {
        CopyFile($"./Externals/Prometheus/prometheus-{prometheusVersion}.windows-amd64/prometheus.exe", $"./bin/{configuration}/prometheus.exe");
        CopyDirectory($"./Externals/Grafana/grafana-{grafanaVersion}", $"./bin/{configuration}/grafana");
    }
});

Task("Publish")
    .IsDependentOn("Clean")
    .IsDependentOn("Build")
    .Does(() =>
{
    // Build dotnet code
    DotNetCorePublish("./Companion/Companion.csproj", new DotNetCorePublishSettings
    {
        Configuration = configuration,
        OutputDirectory = $"./bin/publish/{configuration}",
        PublishSingleFile = true,
        SelfContained = true,
        Runtime = "win-x64",
        MSBuildSettings = msBuildSettings
    });

    DotNetCorePublish("./PrometheusExporter/PrometheusExporter.csproj", new DotNetCorePublishSettings
    {
        Configuration = configuration,
        OutputDirectory = $"./bin/publish/{configuration}",
        PublishSingleFile = true,
        SelfContained = true,
        Runtime = "win-x64",
        MSBuildSettings = msBuildSettings
    });

    CopyFile($"./Externals/Prometheus/prometheus-{prometheusVersion}.windows-amd64/prometheus.exe", $"./bin/publish/{configuration}/prometheus.exe");
    CopyDirectory($"./Externals/Grafana/grafana-{grafanaVersion}", $"./bin/publish/{configuration}/grafana");

    // Build javascript code
    NpmInstall(settings => settings.WorkingDirectory = "./map/");
    NpmRunScript("compile", settings => settings.WorkingDirectory = "./map/");
    CreateDirectory($"bin/publish/{configuration}/map");
    CopyDirectory("./map/vendor", $"bin/publish/{configuration}/map/vendor");
    CopyDirectory("./map/js", $"bin/publish/{configuration}/map/js");
    CopyDirectory("./map/img", $"bin/publish/{configuration}/map/img");
    CopyFile("./map/index.html", $"bin/publish/{configuration}/map/index.html");
    CopyFile("./map/map-16k.png", $"bin/publish/{configuration}/map/map-16k.png");

    // Package
    string version = System.IO.File.ReadAllText("version.txt");
    Zip($"bin/publish/{configuration}", $"./bin/publish/FicsitRemoteMonitoringCompanion-v{version}.zip");
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

enum VersionBumpComponent {
    MAJOR,
    MINOR,
    PATCH
}

var bumpTarget = Argument<string>("bump", "minor");
Task("CutRelease")
    .Does(() => 
{
    string version = System.IO.File.ReadAllText("version.txt");
    var semVersion = SemVersion.Parse(version, false);
    SemVersion nextVer = null;

    try 
    {
        VersionBumpComponent component = (VersionBumpComponent)Enum.Parse(typeof(VersionBumpComponent), bumpTarget.ToUpper());


        switch(component)
        {
            case VersionBumpComponent.MAJOR: 
                nextVer = semVersion.Change(semVersion.Major + 1, 0, 0);
                Console.WriteLine("New version is " + nextVer.ToString());
                break;

            case VersionBumpComponent.MINOR:
                nextVer = semVersion.Change(semVersion.Major, semVersion.Minor + 1, 0);
                Console.WriteLine("New version is " + nextVer.ToString());
                break;

            case VersionBumpComponent.PATCH:
                nextVer = semVersion.Change(semVersion.Major, semVersion.Minor, semVersion.Patch + 1);
                Console.WriteLine("New version is " + nextVer.ToString());
                break;
        }

        System.IO.File.WriteAllText("version.txt", nextVer.ToString());
        GitAdd(".", new FilePath[] { "./version.txt" });
        var taggedCommit = GitCommit(".", "Andy Hunt", "github@andyhunt.me", "Bump version to " + nextVer.ToString());
        GitTag(".", nextVer.ToString());

        Console.WriteLine("Version bumped in commit " + taggedCommit.Sha.Substring(0, 12) + "");
        Console.WriteLine("Run the following to push the new version");
        Console.WriteLine("\tgit push origin main");
        Console.WriteLine("\tgit push origin " + nextVer.ToString());
    }
    catch(ArgumentException)
    {
        Console.WriteLine("bump argument must be one of: major, minor, patch");
    }
}); 

Task("Git-Chglog")
    .IsDependentOn("ExternalsDir")
    .Does(() =>
{
    System.IO.Directory.CreateDirectory("./Externals/Git-Chglog");

    if(!System.IO.File.Exists("./Externals/Git-Chglog/git-chglog.exe"))
        {
        Information($"Downloading git-chglog {gitchglogVersion}");
        DownloadFile(
            $"https://github.com/git-chglog/git-chglog/releases/download/v{gitchglogVersion}/git-chglog_{gitchglogVersion}_windows_amd64.zip",
            "./Externals/Git-Chglog/chglog.zip"
        );


        Unzip("./Externals/Git-Chglog/chglog.zip", "./Externals/Git-Chglog/");
    } 
    else 
    {
       Information("Skipping Git-Chglog download because it has already been downloaded");    
    }
});

Task("ReleaseNotes")
    .IsDependentOn("Git-Chglog")
    .Does(() =>
{
    var chglogPath = "./Externals/Git-Chglog/git-chglog.exe";

    string version = System.IO.File.ReadAllText("version.txt");
    
    string changelogTemplate = System.IO.File.ReadAllText("./.chglog/CHANGELOG.tpl.md");
    string installationInstructions = System.IO.File.ReadAllText("./InstallationInstructions.md");
    string combined = changelogTemplate + Environment.NewLine + Environment.NewLine + installationInstructions;

    System.IO.File.WriteAllText("./Externals/Git-Chglog/template.tpl", combined);

    StartProcess(
      chglogPath, 
      new ProcessSettings {
          Arguments = $"-t ./Externals/Git-Chglog/template.tpl -o ReleaseNotes.md {version}..{version}"
      }
    );
});

//////////////////////////////////////////////////////////////////////
// EXECUTION
//////////////////////////////////////////////////////////////////////

RunTarget(target);