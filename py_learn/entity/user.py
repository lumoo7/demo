from dataclasses import dataclass


@dataclass
class User(object):
    id: int
    code: int
    name: str
    nickname: str

    def print_info(self):
        print("id:{},code:{},name:{},nickname:{}".format(self.id, self.code, self.name, self.nickname))
