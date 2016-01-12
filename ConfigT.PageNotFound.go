package masterbot

func (cfg ConfigT) PageNotFound() []byte {
   return []byte(cfg.PageNotFoundPath)
}

