from ratemyprofessor import School, Professor
import requests
import bs4 as bs
import re

college = School(224)

def query_rmp(professor_name: str) -> bool | list:
    url = 'https://www.ratemyprofessors.com/search/professors/%s?q=%s' % (college.id, professor_name)
    response = requests.get(url)
    soup = bs.BeautifulSoup(response.text, 'html.parser')
    headerDiv: bs.Tag = soup.find('div', {'data-testid': 'search-results-header'})
    queryText = headerDiv.get_text()
    
    if "No professors with" in queryText:
        return False
    else:
        professors_list = []
        profs = re.findall(r'"legacyId":(\d+)', response.text)
        for professor_data in profs:
            try:
                professors_list.append(Professor(int(professor_data)))
            except ValueError:
                pass
        return professors_list

def get_professor_url(professor_name: str) -> str | None:
    if result := query_rmp(professor_name):
        return f"https://www.ratemyprofessors.com/professor/{result[0].id}"
    else:
        return None

if __name__ == "__main__":
    print(get_professor_url("Erik"))
    print(get_professor_url("Abrar"))