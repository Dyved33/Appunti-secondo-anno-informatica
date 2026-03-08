import database
from gui import avvia_gui

def main():
    database.init_database()
    avvia_gui()

if __name__ == "__main__":
    main()
