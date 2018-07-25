import sys
sys.path.append('../')

import cla
from cla.utils import create_database, delete_database
delete_database()
create_database()
