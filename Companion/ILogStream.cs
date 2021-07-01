namespace Companion
{
    internal interface ILogStream
    {
        string FullLogOutput { get; }

        event LogLineArrived OnLogLine;
    }
}