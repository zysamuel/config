RM=rm -f
SRCS=apihandlers.go\
	  ipcutils.go\
	  client.go\
	  clientmap.go\
	  clientif.go\
	  logger.go\
	  apiroutes.go\
	  restroutes.go\
	  configmgr.go\
	  main.go\
	  remotebgppeer.go

COMP_NAME=confd
all: $(COMP_NAME)

$(COMP_NAME): $(SRCS)
	 go build -o $@ $(SRCS)

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

clean:guard
	 $(RM) $(COMP_NAME) 
