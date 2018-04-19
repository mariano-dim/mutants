Resumen solucion (rationale):

Para soportar la carga o mejor dicho, para soportar una carga flexible, se ha creado un cluster de Kubernates en Digital Ocean, formado por tres nodos Master y dos nodos workers, en la zona 2 de San Francisco.

La creacion del cluster con algunas herramientas preconfiguradas se realizo mediante stackpoint.io, que facilita esta tarea (trial por un mes).

Se realizaron algunas pruebas de carga utilizando SOAPUI. El archivo del proyecto se encuentra en este mismo repositorio en la carpeta SoapUIFiles.

Para soportar una carga de 1.000.000 de request por segundo es necesario contar con una infraestructura acorde!!! Y lo mismo para poder enviar requerimientos a ese ratio. En el Test contamos con un cluster de kubernates (que es pago!) y es dificil alcanzar estos requerimientos, pero se demuestra una solucion para el escalado horizontal de la API y del Motor de dase de datos (y vertical).

Acrtualmente en el cluster de kubernates hay 12 replicas del contenedor api rest y tres de mongodb.

Se opto por una solucion de base de datos no sql, mongodb. El motor se instalo en kubernetes mediante un conjunto de contenedores statefullset. Mongodb se ideo de forma nativa para escalar horizontalmente (si aceptamos trabajar con eventual consistency).

Para la api rest se opto por un conjunto de servicios steteless.

Ambos conjuntos pueden escalar horizonaltamente. Kubernetes crea automaticamente con nuestros servicios una red ingress, que basicamente provee un Loadbalancer interno. Este servicio de tipo NodePort, junto con la direccion IP de uno de los nodos master, permiten conectar con los endpoints de la solucion (/stats /mutants). El tipo de servicio NodePort enmascara el puerto 8080 de nuestra api rest. Al setear Servicios NodePort, stackpoint.io configura automaticamente un LoadBalancer sobre cada nodo master de Kubernates. En este caso como tenemos tres masters, hay un LB apuntando a tres Masters

NodePort: 31227/

Ejemplo del servicio mutants: curl -i -X POST -H "Content-Type:application/json" -d'{ "dna" : ["TTAAAT", "TACTCC", "ATACAC", "AAGACT", "CCACTT", "ATGAAT"]}' http://165.227.13.175:31227/mutant

Los otros endpoints disponibles son:

(GET) http://165.227.13.175:31227/ping

(GET) http://165.227.13.175:31227/stats

Ls solucion api rest se desarrollo con golang, se utilizaron los packetes extras:

github.com/gorilla/mux

gopkg.in/mgo.v2

gopkg.in/mgo.v2/bson

Debido a la consigna del problema (HA + Performance) se opto por utilizar estas librerias basicas y no frameworks como revel u otros.

Se analizaron varias opciones para que el algoritmo de descubrimiento sea eficiente, incluyendo algunas soluciones con backtraking y poda, pero se fueron descartando debido a que esta solucion mostro ser mas eficiente (es una solucion de busqueda directa).

Para poder despegar la solucion en cloud, fue necesario publicar la imagen Dockerfile en docker hub. La imagen se buildeo localmente y el Dockerfile se encuentra en el repositorio, lo mismo que los archivos para pruebas locales docker-compose.yml y los archivos para desplegar el deployment y el servicio en kubernates.

Para interactuar con el cluster de kubernates se utiliza el cliente kubectl, a traves de este se han creado todos los objetos relacionados con la solucion.

TODO Para mejorar ls solucion (enpoint stats) se puede incluir un redis para utilizar las funciones del estilo de INCR para que todas las instancias stateless del api rest, puedan mantener elementos compartidos, como lo son las estadisticas de Humans y Mutants. Y tal vez un grafana para mostrar las estadisticas.

Se pueden incluir mejoras a nivel api rest, para evitar algunos tipos de ataques o para controlar la cantidad de requerimientos que puede atender la solucion por unidad de tiempo. Otras importantes mejoras que se deben hacer son: externalizar la configuracion, internacionalizar la api rest, exportar las metricas a un grafana, flexibilizar la solcion de kubernates para que el escalado de contenedores sea elastico.

Se debe documentar la API, por ejemplo con swagger o similar.

Tambien se debe chequear el correcto uso de las conexiones de base de datos y los servicios asociados como asi tambien mejorar un poco el codigo. Por ejemplo desacomplar el acceso a la base de datos de la logica de negocio principal.

Consideraciones acerca de los algoritmos utilizados para identificar secuencias

Una vez que se encuentran cuatro elementos iguales en una fila, columna o diagonal, damos por hecho ese proceso, abandonando la fila, columna o diagonal.

Las secuencias encontradas se van registrando en una matriz para poder identificar cuales ya fueron encontrados porque pertenecian a otra secuencia

El orden en el cual se van hallando las secuencias importa, es decir que si por ejemplo, hay una secuencia de forma horizontal y hay un pivot vertical, posiblemente la secuencia horizonal elimine una posible solucion vertical.

Con respecto a los Test unitarios, comparto el resumen donde se ve que superan el 80% de cobertura

MacBook-Pro-de-Mariano:api-rest marianoandresdimaggio$ go tool cover -func=coverage.out

ml/com/mutants/api-rest/actions.go:14:	getSession	75.0%

ml/com/mutants/api-rest/actions.go:27:	Close	0.0%

ml/com/mutants/api-rest/actions.go:34:	Ping	100.0%

ml/com/mutants/api-rest/actions.go:39:	StatsMutants	73.3%

ml/com/mutants/api-rest/actions.go:74:	MutantCheck	82.1%

ml/com/mutants/api-rest/actions.go:130:	responseAdnCheck	100.0%

ml/com/mutants/api-rest/actions.go:136:	responseStats	100.0%

ml/com/mutants/api-rest/actions.go:143:	isMutant	100.0%

ml/com/mutants/api-rest/actions.go:206:	IsMutantByExternalObliqueAnalysis	75.0%

ml/com/mutants/api-rest/actions.go:269:	IsMutantByLateralObliqueAnalysis	75.0%

ml/com/mutants/api-rest/actions.go:330:	IsMutantByHorizontalAnalysis	89.5%

ml/com/mutants/api-rest/actions.go:368:	IsMutantByVerticalAnalysis	100.0%

ml/com/mutants/api-rest/actions.go:406:	IsMutantByCentralObliqueAnalysis	76.5%

ml/com/mutants/api-rest/actions.go:453:	PrintMatrix	100.0%

ml/com/mutants/api-rest/actions.go:466:	isValidChainElements	100.0%

ml/com/mutants/api-rest/main.go:5:	main	0.0%

ml/com/mutants/api-rest/routes.go:28:	Run	0.0%

ml/com/mutants/api-rest/routes.go:34:	Initialize	0.0%

total:	(statements)	81.4%

listado de PODS de la solucion en el cluster

MacBook-Pro-de-Mariano:mutants marianoandresdimaggio$ kubectl get pods

NAME READY STATUS RESTARTS AGE

api-rest-ml-547998d745-69kds 1/1 Running 0 3h

api-rest-ml-547998d745-6mgft 1/1 Running 0 3h

api-rest-ml-547998d745-6rgzt 1/1 Running 0 3h

api-rest-ml-547998d745-d9ghb 1/1 Running 0 8s

api-rest-ml-547998d745-fz5bx 1/1 Running 0 3h

api-rest-ml-547998d745-ksn62 1/1 Running 0 3h

api-rest-ml-547998d745-lrndq 1/1 Running 0 8s

api-rest-ml-547998d745-pzgfj 1/1 Running 0 8s

api-rest-ml-547998d745-ttsjx 1/1 Running 0 8s

api-rest-ml-547998d745-vbndv 1/1 Running 0 8s

api-rest-ml-547998d745-vvntd 1/1 Running 0 3h

api-rest-ml-547998d745-x8r98 1/1 Running 0 8s

ml-cluster-mongodb-replicaset-0 1/1 Running 0 4h

ml-cluster-mongodb-replicaset-1 1/1 Running 0 4h

ml-cluster-mongodb-replicaset-2 1/1 Running 0 4h

Links de referencia principales

https://cloud.digitalocean.com

https://stackpoint.io

https://hub.docker.com/r/marianodim/api-rest/

https://kubernetes.io/blog/2017/01/running-mongodb-on-kubernetes-with-statefulsets

https://hackernoon.com/build-restful-api-in-go-and-mongodb-5e7f2ec4be94

https://kublr.com/blog/how-to-run-a-mongodb-replica-set-on-kubernetes-petset-or-statefulset/

https://kubernetes.io/docs/reference/kubectl/cheatsheet/
