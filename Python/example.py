import collections  # From Python standard library.
import bson
from bson.codec_options import CodecOptions
import aws-cli 
import chacha20
import boto
import boto3
import botocore
import botocore.session

data = bson.BSON.encode({'a': 1})
decoded_doc = bson.BSON(data).decode()
print(type(decoded_doc)) 


options = CodecOptions(document_class=collections.OrderedDict)
decoded_doc = bson.BSON(data).decode(codec_options=options)
print(type(decoded_doc)) 
