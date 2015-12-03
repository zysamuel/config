RM=rm -f
DESTDIR=$(SR_CODE_BASE)/snaproute/src/bin
PARAMSDIR=$(DESTDIR)/params
MKDIR=mkdir -p

SRCS=apihandlers.go\
	  ipcutils.go\
	  client.go\
	  clientmap.go\
	  clientif.go\
	  objif.go\
	  logger.go\
	  restroutes.go\
	  configmgr.go\
	  dbif.go\
	  main.go\
	  remotebgppeer.go

COMP_NAME=confd
all: exe install 

exe: $(SRCS)
	 go build -o $(DESTDIR)/$(COMP_NAME) $(SRCS)

install:
	 $(MKDIR) $(PARAMSDIR)
	 install params/clients.json $(PARAMSDIR)/
	 install $(SR_CODE_BASE)/snaproute/src/models/objectconfig.json $(PARAMSDIR)

fmt: $(SRCS)
	 go fmt $(SRCS)

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

clean:guard
	 $(RM) $(DESTDIR)/$(COMP_NAME) 
