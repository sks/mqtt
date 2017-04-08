MQTT
====

This is a mqtt_ consumer_ written for Gollum_.

Parameters
----------

**Enable**
 Enable switches the consumer on or off.
 By default this value is set to true.

**ID**
 ID allows this consumer to be found by other plugins by name.
 By default this is set to "" which does not register this consumer.

**Stream**
 Stream contains either a single string or a list of strings defining the message channels this consumer will produce.
 By default this is set to "*" which means only producers set to consume "all streams" will get these messages.

**Fuse**
 Fuse defines the name of a fuse to observe for this consumer.
 Producer may "burn" the fuse when they encounter errors.
 Consumers may react on this by e.g. closing connections to notify any writing services of the problem.
 Set to "" by default which disables the fuse feature for this consumer.
 It is up to the consumer implementation to react on a broken fuse in an appropriate manner.

**connectionString**
 MQTT Connection string.
 This denotes the MQTT Server Connection string.
 By Default, this starts listening to *tcp://localhost:1883*


**topic**
  MQTT Topic to start listening to. By Default this listens to *#*

Sample config
-------------
.. code-block:: yaml

    "Mqtt":
        Type: "mqtt.MqttConsumer"
        Streams: "Mqtt"
        topic: testtopic
        Fuse: "Mqtt"

    "FileOut":
        Type: "producer.File"
        Streams: "Mqtt"
        Fuse: "Mqtt"
        Modulators:
            - format.Envelope:
                Postfix: "\n"
        File: /tmp/gollum_test.log
        Batch:
            TimeoutSec: 1



.. [MQTT] http://mqtt.org/
.. [consumer] http://gollum.readthedocs.io/en/latest/consumers/index.html)
.. [Gollum] https://github.com/trivago/gollum
