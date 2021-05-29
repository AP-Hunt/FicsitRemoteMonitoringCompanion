using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;

namespace Companion.Config
{
    /// <summary>
    /// Interaction logic for ConfigWindow.xaml
    /// </summary>
    public partial class ConfigWindow : Window
    {
        public ConfigFile Config = null;

        public ConfigWindow()
        {
            InitializeComponent();

            Task<ConfigFile> t = ConfigIO.ReadConfigFile(Paths.ConfigFile);
            t.Wait();
            this.Config = t.Result;
            this.txtSatisfactoryGameDir.Text = this.Config.SatisfactoryGameDirectory;
        }

        private void btnFindGameDir_Click(object sender, RoutedEventArgs e)
        {
            string currentPath = this.txtSatisfactoryGameDir.Text;
            System.Windows.Forms.FolderBrowserDialog fbd = new System.Windows.Forms.FolderBrowserDialog();

            if (currentPath != "")
            {
                fbd.SelectedPath = currentPath;
            }
            
            if(fbd.ShowDialog() == System.Windows.Forms.DialogResult.OK)
            {
                this.txtSatisfactoryGameDir.Text = fbd.SelectedPath;
            }
        }

        private async void btnOK_Click(object sender, RoutedEventArgs e)
        {
            ConfigFile cfg = await ConfigIO.ReadConfigFile(Paths.ConfigFile);
            cfg.SatisfactoryGameDirectory = this.txtSatisfactoryGameDir.Text;

            if(!cfg.IsValidGamePath())
            {
                MessageBox.Show("Invalid Satsifactory game path. Could not find FactoryGame.exe.");
                return;
            }

            await ConfigIO.WriteConfigFile(Paths.ConfigFile, cfg);
            this.Config = cfg;
            this.Close();
        }
    }
}
