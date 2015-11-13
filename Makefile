RM=rm -f
DESTDIR=$(SR_CODE_BASE)/snaproute/src/bin
SRCS=apihandlers.go\
	  ipcutils.go\
	  client.go\
	  clientmap.go\
	  clientif.go\
	  objif.go\
	  logger.go\
	  restroutes.go\
	  configmgr.go\
	  main.go\
	  remotebgppeer.go

COMP_NAME=confd
all: exe

exe: $(SRCS)
	 go build -o $(DESTDIR)/$(COMP_NAME) $(SRCS)

fmt: $(SRCS)
	 go fmt $(SRCS)

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

clean:guard
	 $(RM) $(COMP_NAME) 
