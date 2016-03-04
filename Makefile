RM=rm -f
DESTDIR=$(SR_CODE_BASE)/snaproute/src/out/bin
PARAMSDIR=$(DESTDIR)/params
MKDIR=mkdir -p
RSYNC=rsync -rupE

CLIENTIF_FILES = $(shell ls  *clientif.go)
SRCS=$(CLIENTIF_FILES)\
	  apihandlers.go\
	  apierrcodes.go\
	  ipcutils.go\
	  client.go\
	  clientmap.go\
	  objif.go\
	  logger.go\
	  restroutes.go\
	  configmgr.go\
	  dbif.go\
	  ipblockmgr.go\
	  usermgmt.go\
	  main.go

COMP_NAME=confd
all: gencode exe install 

exe: $(SRCS)
	 go build -o $(DESTDIR)/$(COMP_NAME) $(SRCS)

install:
	 @$(MKDIR) $(PARAMSDIR)
	 @$(RSYNC) docsui $(PARAMSDIR)
	 @echo $(DESTDIR)
	 install params/* $(PARAMSDIR)/
	 install $(SR_CODE_BASE)/snaproute/src/models/objectconfig.json $(PARAMSDIR)
	 install $(SR_CODE_BASE)/snaproute/src/models/genObjectConfig.json $(PARAMSDIR)

fmt: $(SRCS)
	 go fmt $(SRCS)

gencode:
	$(SR_CODE_BASE)/reltools/codegentools/gencode.sh

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

clean:guard
	 $(RM) $(DESTDIR)/$(COMP_NAME) 
