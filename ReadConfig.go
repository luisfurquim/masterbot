package masterbot

import (
   "os"
   "io"
   "io/ioutil"
   "encoding/json"
)

func ReadConfig(cfgReader io.Reader) ([]byte, *ConfigT, error) {
   var ConfigFile []byte
   var Config      *ConfigT
   var err          error

   ConfigFile, err = ioutil.ReadAll(cfgReader)
   if err != nil {
      Goose.Logf(1,"%s (%s)",ErrReadingConfig,err)
      return nil, nil, ErrReadingConfig
   }

   Config = &ConfigT{}
   err = json.Unmarshal(ConfigFile,Config)
   if err != nil {
      Goose.Logf(1,"%s (%s) %s",ErrParsingConfig,err,ConfigFile)
      return nil, nil, ErrParsingConfig
   }

   os.Chdir(Config.WorkDir)

   err = Config.LoadClientCert()
   if err!=nil {
      Goose.Logf(1,"%s (%s)",ErrLoadingCliCerts,err)
      return nil, nil, ErrLoadingCliCerts
   }

   return ConfigFile, Config, nil
}

