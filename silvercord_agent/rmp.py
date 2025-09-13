from ratemyprofessor import Professor
import requests
import bs4 as bs
import re
from constants import CCNY_SCHOOL_ID

headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
}

def query_rmp(college_id: str, professor_name: str) -> (list, bool):
    url = 'https://www.ratemyprofessors.com/search/professors/%s?q=%s' % (college_id, professor_name)
    response = requests.get(url, headers=headers)
    soup = bs.BeautifulSoup(response.text, 'html.parser')
    headerDiv: bs.Tag = soup.find('div', {'data-testid': 'search-results-header'})
    queryText = headerDiv.get_text() # NOTE: This finds div header (e.g. "4 professors with 'ProfName' in their name at 'CollegeName'")

    if "No professors with" in queryText:
        return [], False
    professors_list = []
    profs = re.findall(r'"legacyId":(\d+)', response.text)
    for professor_data in profs:
        try:
            professors_list.append(Professor(int(professor_data)))
        except ValueError:
            pass
    return professors_list, True

def get_professor_url(professor_name: str) -> str | None:
    if result := query_rmp(CCNY_SCHOOL_ID, professor_name):
        return f"https://www.ratemyprofessors.com/professor/{result[0].id}"
    else:
        return None