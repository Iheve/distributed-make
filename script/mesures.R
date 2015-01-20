
data= read.csv("mesures.csv", sep=";")

library(ggplot2);

for (i in 1:length(data$nthreads)) {
    if (data[i, "nthreads"] == 1) {
        tempsseq= data[i,"temps"]
    }
}

d = data.frame(nthreads=data$nthreads, temps=data$temps);

p = ggplot(d, aes(y=temps,x=nthreads)) + geom_line()
p
ggsave("temps.svg",width=2*par("din")[1])

p = ggplot(d, aes(y=temps*nthreads,x=nthreads)) + geom_line()
p
ggsave("travail.svg",width=2*par("din")[1])

p = ggplot(d, aes(y=tempsseq/temps,x=nthreads)) + geom_line()
p
ggsave("acceleration.svg",width=2*par("din")[1])

p = ggplot(d, aes(y=tempsseq/(temps*nthreads),x=nthreads)) + geom_line()
p
ggsave("efficacite.svg",width=2*par("din")[1])
