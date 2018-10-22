/*
 * Ceph Nano (C) 2017 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cliKubeNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube",
		Short: "Outputs cn kubernetes template (cn kube > kube-cn.yml)",
		Args:  cobra.NoArgs,
		Run:   kubeTemplate,
	}

	return cmd
}

func kubeTemplate(cmd *cobra.Command, args []string) {
	kubeTmp := `
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cn-pv-claim-var
  labels:
    app: ceph
	  daemon: nano
spec:
  # Read more about access modes here: http://kubernetes.io/docs/user-guide/persistent-volumes/#access-modes
  accessModes:
	- ReadWriteOnce
  resources:
	requests:
      storage: 10Gi
  # Uncomment and add storageClass specific to your requirements below. Read more https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class-1
  #storageClassName:
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cn-pv-claim-etc
  labels:
	app: ceph
	  daemon: nano
spec:
  # Read more about access modes here: http://kubernetes.io/docs/user-guide/persistent-volumes/#access-modes
  accessModes:
	- ReadWriteOnce
  resources:
	requests:
	storage: 10Mi
  # Uncomment and add storageClass specific to your requirements below. Read more https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class-1
  #storageClassName:
---
  apiVersion: v1
  kind: Service
  metadata:
    name: ceph-nano-s3
    labels:
      app: ceph
      daemon: nano
  spec:
    ports:
    - name: cn-s3
      port: 80
      protocol: TCP
      targetPort: 8000
    type: LoadBalancer
    selector:
      app: ceph
      daemon: nano
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: ceph
    daemon: nano
    name: ceph-nano
  spec:
    replicas: 1
    serviceName: ceph-nano
    selector:
      matchLabels:
        app: ceph
    template:
      metadata:
        name: ceph-nano
        labels:
  	      app: ceph
  	      daemon: nano
  	  spec:
  	    containers:
  	    - image: ceph/daemon
  	      imagePullPolicy: Always
  	      name: ceph-nano
  	      ports:
  	      - containerPort: 8000
  	    	name: cn-s3
  	    	protocol: TCP
  	    	resources:
  	    	  limits:
  	    	    cpu: "1"
  	    	    memory: 512M
  	    	  requests:
  	    	    cpu: "1"
  	    	    memory: 512M
  	    	env:
  	    	- name: NETWORK_AUTO_DETECT
  	    	  value: "4"
  	    	- name: RGW_CIVETWEB_PORT
  	    	  value: "8000"
  	    	- name: SREE_PORT
  	    	  value: "5001"
  	    	- name: CEPH_DEMO_UID
  	    	  value: "nano"
  	    	- name: CEPH_DAEMON
  	    	  value: "demo"
  	    	- name: DEBUG
  	    	  value: "verbose"
            volumeMounts:
            - name: cn-varlibceph
              mountPath: /var/lib/ceph
            - name: cn-etcceph
              mountPath: /etc/ceph
			volumes:
			- name: cn-varlibceph
			  persistentVolumeClaim:
				 claimName: cn-pv-claim-var
			- name: cn-etcceph
			  persistentVolumeClaim:
					claimName: cn-pv-claim-etc
`
	fmt.Print(kubeTmp)
}
