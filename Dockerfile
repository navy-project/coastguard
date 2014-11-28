FROM scratch

ADD bin/coastguard /coastguard

CMD ["/coastguard"]
