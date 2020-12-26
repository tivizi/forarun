FROM ubuntu:20.10
RUN apt update;apt install ca-certificates -y
ADD bin/forarun /
RUN mkdir /config
ADD approot/config/app.yaml /config/app.yaml 
RUN mkdir /templates
ADD approot/templates /templates
ADD approot/assets /assets
ADD approot/ip.dat /ip.dat
ENV TZ Asia/Shanghai
CMD [ "/forarun" ]
